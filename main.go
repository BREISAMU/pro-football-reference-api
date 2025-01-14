package main

import (
	"log"
	"net/http"
	handlers "pfr/handlers"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
Gets general history overlook, see "https://www.pro-football-reference.com/teams/" links
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getSeasonOverlook(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	year_int, err_atoi := strconv.Atoi(year)

	if err_atoi != nil {
		log.Println("Error parsing specified year.")
		return
	}

	url := "https://www.pro-football-reference.com/teams/" + team + "/"
	tableSelector := "#team_index"

	data, err := handlers.GetSeasonOverlook(url, tableSelector, year_int)

	if err != nil {
		log.Println("Error retrieving team history data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

/*
Gets summary of specific draft year by team, see "https://www.pro-football-reference.com/teams/buf/draft.htm" as example with param "buf"
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getDraftYear(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	year_int, err_atoi := strconv.Atoi(year)

	if err_atoi != nil {
		log.Println("Error parsing specified year.")
		return
	}

	url := "https://www.pro-football-reference.com/teams/" + team + "/draft.htm"
	tableSelector := "#draft"

	data, err := handlers.GetDraftYear(url, tableSelector, year_int)

	if err != nil {
		log.Println("Error retrieving draft data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

/*
Gets offensive stats by team and year, see "https://www.pro-football-reference.com/teams/rav/2024.htm" first table as example with param "rav"
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getTeamOffensiveStats(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	url := "https://www.pro-football-reference.com/teams/" + team + "/" + year + ".htm"
	tableSelector := "#team_stats"

	data, _, _, _, err := handlers.GetTeamYearStats(url, tableSelector, year, team)

	if err != nil {
		log.Println("Error retrieving draft data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

/*
Gets defensive stats by team and year, see "https://www.pro-football-reference.com/teams/rav/2024.htm" first table as example with param "rav"
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getTeamDefensiveStats(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	url := "https://www.pro-football-reference.com/teams/" + team + "/" + year + ".htm"
	tableSelector := "#team_stats"

	_, data, _, _, err := handlers.GetTeamYearStats(url, tableSelector, year, team)

	if err != nil {
		log.Println("Error retrieving draft data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

/*
Gets offensive rankings by team and year, see "https://www.pro-football-reference.com/teams/rav/2024.htm" first table as example with param "rav"
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getTeamOffensiveRankings(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	url := "https://www.pro-football-reference.com/teams/" + team + "/" + year + ".htm"
	tableSelector := "#team_stats"

	_, _, data, _, err := handlers.GetTeamYearStats(url, tableSelector, year, team)

	if err != nil {
		log.Println("Error retrieving draft data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

/*
Gets defensive rankings by team and year, see "https://www.pro-football-reference.com/teams/rav/2024.htm" first table as example with param "rav"
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getTeamDefensiveRankings(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	url := "https://www.pro-football-reference.com/teams/" + team + "/" + year + ".htm"
	tableSelector := "#team_stats"

	_, _, _, data, err := handlers.GetTeamYearStats(url, tableSelector, year, team)

	if err != nil {
		log.Println("Error retrieving draft data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func main() {
	router := gin.Default()

	router.GET("/team/", getSeasonOverlook)                         // ?team=___&year=___
	router.GET("/team/draft", getDraftYear)                         // ?team=___&year=___
	router.GET("/team/offensiveStats", getTeamOffensiveStats)       // ?team=___&year=___
	router.GET("/team/defensiveStats", getTeamDefensiveStats)       // ?team=___&year=___
	router.GET("/team/offensiveRankings", getTeamOffensiveRankings) // ?team=___&year=___
	router.GET("/team/defensiveRankings", getTeamDefensiveRankings) // ?team=___&year=___

	router.Run()
}
