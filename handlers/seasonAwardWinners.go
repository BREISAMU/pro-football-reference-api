package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type AwardWinner struct {
	Award  string `json:"award"`
	Winner string `json:"winner"`
}

func GetSeasonAwardWinners(url string) ([]AwardWinner, error) {
	// Begin client boilerplate
	client := &http.Client{
		Timeout: 32 * time.Second,
	}
	maxRetries := 2
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []AwardWinner{}, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers to mimic a browser
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return []AwardWinner{}, fmt.Errorf("error making request: %v", err)
		}

		// Handle rate limit (429)
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return []AwardWinner{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 10 * time.Second // Default wait time
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(seconds) * time.Second
				}
			}

			log.Printf("Rate limited. Waiting %v before retry %d/%d", waitTime, attempt, maxRetries)
			time.Sleep(waitTime)
			continue
		}

		// Success case
		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	// End client boilerplate

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []AwardWinner{}, fmt.Errorf("error parsing doc: %v", err)
	}

	// Take placeholder data because real is loaded dynamically
	// *appears to always be up to date
	divSelection := doc.Find("#all_awards")
	html, err := divSelection.Html()
	if err != nil {
		return []AwardWinner{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	// Clean html, remove comment symbols and ":"
	commentMarkersRegex := regexp.MustCompile(`<!--|-->|:`)
	cleanHtml := commentMarkersRegex.ReplaceAllString(html, "")

	// Create new doc from parsed and cleaned HTML comment
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(cleanHtml))
	if err != nil {
		return []AwardWinner{}, fmt.Errorf("error cleaning HTML: %v", err)
	}
	divSelection = doc.Find("#div_awards")

	// Setup result values
	var awardWinners []AwardWinner
	var awardWinnerHolder AwardWinner

	// Extract awards and winners
	isCategory := true
	divSelection.Contents().Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) != "" {
			if isCategory {
				awardWinnerHolder.Award = s.Text()
			} else {
				awardWinnerHolder.Winner = s.Text()
				awardWinners = append(awardWinners, awardWinnerHolder)
			}
			isCategory = !isCategory
		}
	})

	if len(awardWinners) == 0 {
		return []AwardWinner{}, fmt.Errorf("no award winners found")
	}

	return awardWinners, nil
}
