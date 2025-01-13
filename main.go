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

func main() {
	router := gin.Default()

	router.GET("/team", getSeasonOverlook) // ?team=___&year=___
	router.GET("/draft", getDraftYear)     // ?team=___&year=___

	router.Run()
}
