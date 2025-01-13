package teamhistory

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SeasonOverlook struct {
	Year               int     `json:"year"`
	League             string  `json:"league"`
	Team               string  `json:"team"`
	Wins               int     `json:"wins"`
	Losses             int     `json:"losses"`
	Ties               int     `json:"ties"`
	DivisionFinish     int     `json:"divisionFinish"`
	PlayoffExitRound   int     `json:"playoffExitRound"`
	PointsFor          int     `json:"pointsFor"`
	PointsAgainst      int     `json:"pointsAgainst"`
	PointsDif          int     `json:"pointsDif"`
	HeadCoaches        string  `json:"headCoaches"`
	BestPlayerAv       string  `json:"bestPlayerAv"`
	BestPlayerPasser   string  `json:"bestPlayerPasser"`
	BestPlayerRusher   string  `json:"bestPlayerRusher"`
	BestPlayerReceiver string  `json:"bestPlayerReceiver"`
	OffRankPts         int     `json:"offRankPts"`
	OffRankYds         int     `json:"offRankYds"`
	DefRankPts         int     `json:"defRankPts"`
	DefRankYds         int     `json:"defRankYds"`
	TakeawayRank       int     `json:"takeawayRank"`
	PointsDifRank      int     `json:"pointsDifRank"`
	YardsDifRank       int     `json:"yardsDifRank"`
	TeamsInLeague      int     `json:"teamsInLeague"`
	MarginOfVictory    float64 `json:"marginOfVictory"`
	StrengthOfSchedule float64 `json:"strengthOfSchedule"`
	Srs                float64 `json:"srs"`
	OffensiveSrs       float64 `json:"offensiveSrs"`
	DefensiveSrs       float64 `json:"defensiveSrs"`
}

func GetTeamHistory(url string, tableSelector string, year int) (SeasonOverlook, error) {
	// Create a custom client with reasonable timeouts
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Maximum number of retries
	maxRetries := 3
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Create a new request with headers
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return SeasonOverlook{}, fmt.Errorf("error creating request: %v", err)
		}

		// Add headers to appear more like a browser
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return SeasonOverlook{}, fmt.Errorf("error making request: %v", err)
		}

		// Check if we hit the rate limit
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return SeasonOverlook{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			// Get retry delay from header or use default
			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 60 * time.Second // default wait time
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(seconds) * time.Second
				}
			}

			log.Printf("Rate limited. Waiting %v before retry %d/%d", waitTime, attempt, maxRetries)
			time.Sleep(waitTime)
			continue
		}

		// Break the retry loop if we got a good response
		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()
		return SeasonOverlook{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// Rest of your existing scraping logic here...
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return SeasonOverlook{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	var tableData [][]string
	var season []string
	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})
		if len(rowData) > 0 {
			tableData = append(tableData, rowData)
			rowYear, err := strconv.Atoi(rowData[0])
			if err == nil && rowYear == year {
				season = rowData
			}
		}
	})

	if len(season) == 0 {
		return SeasonOverlook{}, fmt.Errorf("no data found for year %d", year)
	}
	if len(season) == 0 {
		return SeasonOverlook{}, err
	}

	league := season[1]
	team := season[2]
	wins, _ := strconv.Atoi(season[3])
	losses, _ := strconv.Atoi(season[4])
	ties, _ := strconv.Atoi(season[5])
	divisionFinish, _ := strconv.Atoi(season[6][0:1])

	playoffExitRoundString := season[7]
	playoffExitRoundInt := 0
	if playoffExitRoundString == "Lost WC" {
		playoffExitRoundInt = 1
	} else if playoffExitRoundString == "Lost Div" {
		playoffExitRoundInt = 2
	} else if playoffExitRoundString == "Lost Conf" {
		playoffExitRoundInt = 3
	} else if playoffExitRoundString == "Lost SB" || playoffExitRoundString == "Lost Champ" {
		playoffExitRoundInt = 4
	} else if playoffExitRoundString == "Won SB" || playoffExitRoundString == "Won Champ" {
		playoffExitRoundInt = 5
	}

	pointsFor, _ := strconv.Atoi(season[8])
	pointsAgainst, _ := strconv.Atoi(season[9])
	pointsDif, _ := strconv.Atoi(season[10])
	headCoaches := season[11]
	bestPlayerAv := season[12]
	bestPlayerPasser := season[13]
	bestPlayerRusher := season[14]
	bestPlayerReceiver := season[15]
	offRankPts, _ := strconv.Atoi(season[16])
	offRankYds, _ := strconv.Atoi(season[17])
	defRankPts, _ := strconv.Atoi(season[18])
	defRankYds, _ := strconv.Atoi(season[19])
	takeawayRank, _ := strconv.Atoi(season[20])
	pointsDifRank, _ := strconv.Atoi(season[21])
	yardsDifRank, _ := strconv.Atoi(season[22])
	teamsInLeague, _ := strconv.Atoi(season[23])
	marginOfVictory, _ := strconv.ParseFloat(season[24], 64)
	strengthOfSchedule, _ := strconv.ParseFloat(season[25], 64)
	srs, _ := strconv.ParseFloat(season[26], 64)
	offensiveSrs, _ := strconv.ParseFloat(season[27], 64)
	defensiveSrs, _ := strconv.ParseFloat(season[28], 64)

	seasonOverlook := SeasonOverlook{
		Year:               year,
		League:             league,
		Team:               team,
		Wins:               wins,
		Losses:             losses,
		Ties:               ties,
		DivisionFinish:     divisionFinish,
		PlayoffExitRound:   playoffExitRoundInt,
		PointsFor:          pointsFor,
		PointsAgainst:      pointsAgainst,
		PointsDif:          pointsDif,
		HeadCoaches:        headCoaches,
		BestPlayerAv:       bestPlayerAv,
		BestPlayerPasser:   bestPlayerPasser,
		BestPlayerRusher:   bestPlayerRusher,
		BestPlayerReceiver: bestPlayerReceiver,
		OffRankPts:         offRankPts,
		OffRankYds:         offRankYds,
		DefRankPts:         defRankPts,
		DefRankYds:         defRankYds,
		TakeawayRank:       takeawayRank,
		PointsDifRank:      pointsDifRank,
		YardsDifRank:       yardsDifRank,
		TeamsInLeague:      teamsInLeague,
		MarginOfVictory:    marginOfVictory,
		StrengthOfSchedule: strengthOfSchedule,
		Srs:                srs,
		OffensiveSrs:       offensiveSrs,
		DefensiveSrs:       defensiveSrs,
	}

	return seasonOverlook, nil
}
