package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TeamSeason struct {
	Team               string  `json:"team"`
	Wins               int     `json:"wins"`
	Losses             int     `json:"losses"`
	WinLossPerc        float64 `json:"winLossPerc"`
	PointsFor          int     `json:"pointsFor"`
	PointsAgainst      int     `json:"pointsAgainst"`
	PointsDif          int     `json:"pointsDif"`
	MarginOfVictory    float64 `json:"marginOfVictory"`
	StrengthOfSchedule float64 `json:"strengthOfSchedule"`
	Srs                float64 `json:"srs"`
	OffensiveSrs       float64 `json:"offensiveSrs"`
	DefensiveSrs       float64 `json:"defensiveSrs"`
}

type Conference struct {
	Name      string     `json:"name"`
	Divisions []Division `json:"divisions"`
}

type Division struct {
	Name  string       `json:"name"`
	Teams []TeamSeason `json:"teams"`
}

func GetLeagueStandingsByYearPre1970(url string) ([]Conference, error) {
	tableSelector := "#NFL"

	// ---- CLIENT BOILERPLATE ----
	client := &http.Client{
		Timeout: 4 * (time.Second + 8),
	}

	maxRetries := 2
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Conference{}, fmt.Errorf("error creating request: %v", err)
		}

		// Headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return []Conference{}, fmt.Errorf("error making request: %v", err)
		}

		// Rate limit check
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return []Conference{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 15 * time.Second
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(seconds) * time.Second
				}
			}

			log.Printf("Rate limited. Waiting %v before retry %d/%d", waitTime, attempt, maxRetries)
			time.Sleep(waitTime)
			continue
		}

		// Successful response
		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()
		return []Conference{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Conference{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	var tableData [][]string
	nfl := Conference{"NFL", []Division{}}
	divisionOfInterest := Division{"No", []TeamSeason{}}

	// ---- END ----

	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})

		if len(rowData) == 1 {
			if divisionOfInterest.Name != "No" {
				nfl.Divisions = append(nfl.Divisions, divisionOfInterest)
			}
			divisionOfInterest = Division{rowData[0], []TeamSeason{}}
		}

		if len(rowData) > 1 && rowData[0] != "Tm" {
			println(rowData[0])
			tableData = append(tableData, rowData)
			if err == nil {

				team := rowData[0]
				if team[len(team)-1] == '*' || team[len(team)-1] == '+' {
					team = team[:len(team)-1]
				}

				wins, _ := strconv.Atoi(rowData[1])
				losses, _ := strconv.Atoi(rowData[2])
				winLossPerc, _ := strconv.ParseFloat(rowData[3], 64)
				pointsFor, _ := strconv.Atoi(rowData[4])
				pointsAgainst, _ := strconv.Atoi(rowData[5])
				pointsDif, _ := strconv.Atoi(rowData[6])
				marginOfVictory, _ := strconv.ParseFloat(rowData[7], 64)
				strengthOfSchedule, _ := strconv.ParseFloat(rowData[8], 64)
				srs, _ := strconv.ParseFloat(rowData[9], 64)
				offensiveSrs, _ := strconv.ParseFloat(rowData[10], 64)
				defensiveSrs, _ := strconv.ParseFloat(rowData[11], 64)

				season := TeamSeason{
					Team:               team,
					Wins:               wins,
					Losses:             losses,
					WinLossPerc:        winLossPerc,
					PointsFor:          pointsFor,
					PointsAgainst:      pointsAgainst,
					PointsDif:          pointsDif,
					MarginOfVictory:    marginOfVictory,
					StrengthOfSchedule: strengthOfSchedule,
					Srs:                srs,
					OffensiveSrs:       offensiveSrs,
					DefensiveSrs:       defensiveSrs,
				}

				divisionOfInterest.Teams = append(divisionOfInterest.Teams, season)
			}
		}
	})

	// catch hanging division
	nfl.Divisions = append(nfl.Divisions, divisionOfInterest)

	var league []Conference
	league = append(league, nfl)

	if len(tableData) < 1 {
		return []Conference{}, fmt.Errorf("no data found for selected year")
	}

	return league, nil
}

func GetLeagueStandingsByYearPost1970(url string) ([]Conference, error) {

	// ---- CLIENT BOILERPLATE ----
	client := &http.Client{
		Timeout: 4 * (time.Second + 8),
	}

	maxRetries := 2
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Conference{}, fmt.Errorf("error creating request: %v", err)
		}

		// Headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return []Conference{}, fmt.Errorf("error making request: %v", err)
		}

		// Rate limit check
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return []Conference{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 15 * time.Second
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(seconds) * time.Second
				}
			}

			log.Printf("Rate limited. Waiting %v before retry %d/%d", waitTime, attempt, maxRetries)
			time.Sleep(waitTime)
			continue
		}

		// Successful response
		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()
		return []Conference{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []Conference{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	// ---- END ----

	var tableData [][]string
	afc := Conference{"AFC", []Division{}}
	nfc := Conference{"NFC", []Division{}}
	divisionOfInterest := Division{"No", []TeamSeason{}}

	// Get AFC tables
	doc.Find("#AFC").Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})

		if len(rowData) == 1 {
			if divisionOfInterest.Name != "No" {
				afc.Divisions = append(afc.Divisions, divisionOfInterest)
			}
			divisionOfInterest = Division{rowData[0], []TeamSeason{}}
		}

		if len(rowData) > 1 && rowData[0] != "Tm" {
			tableData = append(tableData, rowData)
			if err == nil {

				team := rowData[0]
				if team[len(team)-1] == '*' || team[len(team)-1] == '+' {
					team = team[:len(team)-1]
				}

				wins, _ := strconv.Atoi(rowData[1])
				losses, _ := strconv.Atoi(rowData[2])
				winLossPerc, _ := strconv.ParseFloat(rowData[3], 64)
				pointsFor, _ := strconv.Atoi(rowData[4])
				pointsAgainst, _ := strconv.Atoi(rowData[5])
				pointsDif, _ := strconv.Atoi(rowData[6])
				marginOfVictory, _ := strconv.ParseFloat(rowData[7], 64)
				strengthOfSchedule, _ := strconv.ParseFloat(rowData[8], 64)
				srs, _ := strconv.ParseFloat(rowData[9], 64)
				offensiveSrs, _ := strconv.ParseFloat(rowData[10], 64)
				defensiveSrs, _ := strconv.ParseFloat(rowData[11], 64)

				season := TeamSeason{
					Team:               team,
					Wins:               wins,
					Losses:             losses,
					WinLossPerc:        winLossPerc,
					PointsFor:          pointsFor,
					PointsAgainst:      pointsAgainst,
					PointsDif:          pointsDif,
					MarginOfVictory:    marginOfVictory,
					StrengthOfSchedule: strengthOfSchedule,
					Srs:                srs,
					OffensiveSrs:       offensiveSrs,
					DefensiveSrs:       defensiveSrs,
				}

				divisionOfInterest.Teams = append(divisionOfInterest.Teams, season)
			}
		}
	})

	afc.Divisions = append(afc.Divisions, divisionOfInterest)

	divisionOfInterest.Name = "No"

	// Get NFC tables
	doc.Find("#NFC").Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})

		if len(rowData) == 1 {
			if divisionOfInterest.Name != "No" {
				nfc.Divisions = append(nfc.Divisions, divisionOfInterest)
			}
			divisionOfInterest = Division{rowData[0], []TeamSeason{}}
		}

		if len(rowData) > 1 && rowData[0] != "Tm" {
			tableData = append(tableData, rowData)
			if err == nil {

				team := rowData[0]
				if team[len(team)-1] == '*' || team[len(team)-1] == '+' {
					team = team[:len(team)-1]
				}

				wins, _ := strconv.Atoi(rowData[1])
				losses, _ := strconv.Atoi(rowData[2])
				winLossPerc, _ := strconv.ParseFloat(rowData[3], 64)
				pointsFor, _ := strconv.Atoi(rowData[4])
				pointsAgainst, _ := strconv.Atoi(rowData[5])
				pointsDif, _ := strconv.Atoi(rowData[6])
				marginOfVictory, _ := strconv.ParseFloat(rowData[7], 64)
				strengthOfSchedule, _ := strconv.ParseFloat(rowData[8], 64)
				srs, _ := strconv.ParseFloat(rowData[9], 64)
				offensiveSrs, _ := strconv.ParseFloat(rowData[10], 64)
				defensiveSrs, _ := strconv.ParseFloat(rowData[11], 64)

				season := TeamSeason{
					Team:               team,
					Wins:               wins,
					Losses:             losses,
					WinLossPerc:        winLossPerc,
					PointsFor:          pointsFor,
					PointsAgainst:      pointsAgainst,
					PointsDif:          pointsDif,
					MarginOfVictory:    marginOfVictory,
					StrengthOfSchedule: strengthOfSchedule,
					Srs:                srs,
					OffensiveSrs:       offensiveSrs,
					DefensiveSrs:       defensiveSrs,
				}

				divisionOfInterest.Teams = append(divisionOfInterest.Teams, season)
			}
		}
	})

	// Get hanging team
	nfc.Divisions = append(nfc.Divisions, divisionOfInterest)

	var league []Conference
	league = append(league, nfc)
	league = append(league, afc)

	if len(league[0].Divisions) < 1 {
		return []Conference{}, fmt.Errorf("no data found for selected year")
	}

	return league, nil
}
