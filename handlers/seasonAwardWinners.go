package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	client := &http.Client{
		Timeout: 32 * time.Second, // Fixed timeout calculation
	}

	maxRetries := 2
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers to mimic a browser
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		// Handle rate limit (429)
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return nil, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}
	// bodyBytes, _ := io.ReadAll(resp.Body)
	// rawHtml := string(bodyBytes)
	// commentMarkersRegex := regexp.MustCompile(`<!--|-->`)
	// cleanHtml := commentMarkersRegex.ReplaceAllString(rawHtml, "")
	// err = os.WriteFile("output.html", []byte(cleanHtml), 0644)

	////

	divSelection := doc.Find("#all_awards")

	html, err := divSelection.Html()
	err = os.WriteFile("output.html", []byte(html), 0644)
	// divSelection.SetHtml("Hello World")
	// html, err = divSelection.Html()

	commentMarkersRegex := regexp.MustCompile(`<!--|-->`)
	cleanHtml := commentMarkersRegex.ReplaceAllString(html, "")

	err = os.WriteFile("outputCLEAN.html", []byte(cleanHtml), 0644)

	////

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(cleanHtml))
	divSelection = doc.Find("#div_awards")
	var awardWinners []AwardWinner

	// Extract awards and winners
	fmt.Println(divSelection.Contents().Length())
	divSelection.Contents().Each(func(i int, s *goquery.Selection) {

		fmt.Printf("Reaching section: %s\n", s.Text())
	})

	if len(awardWinners) == 0 {
		return nil, fmt.Errorf("no award winners found")
	}

	return awardWinners, nil
}
