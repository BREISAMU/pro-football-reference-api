package main

import (
	"log"
	"net/http"
	teamhistory "pfr/getFunctions"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
Gets general history overlook, see "https://www.pro-football-reference.com/teams/" links
Specify:
- team (gnb, dal, jax, etc.)
- season (2003, 2024, etc.)
*/
func getTeamHistory(c *gin.Context) {
	team := c.Query("team")
	year := c.Query("year")
	year_int, err_atoi := strconv.Atoi(year)

	if err_atoi != nil {
		log.Println("Error parsing specified year.")
		return
	}

	url := "https://www.pro-football-reference.com/teams/" + team + "/"
	tableSelector := "#team_index"

	data, err := teamhistory.GetTeamHistory(url, tableSelector, year_int)

	if err != nil {
		log.Println("Error retrieving team history data.")
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func main() {
	router := gin.Default()

	router.GET("/teamHistory", getTeamHistory)

	router.Run("localhost:8080")
}
