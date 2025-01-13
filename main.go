package main

import (
	"log"
	"net/http"
	teamhistory "pfr/getFunctions"

	"github.com/gin-gonic/gin"
)

func getTeamHistory(c *gin.Context) {
	// team := c.Param("team")
	team := "gnb"
	// year, err_atoi := strconv.Atoi(c.Param("year"))

	year := 2000
	// if err_atoi != nil {
	// 	log.Println("Error parsing specified year.")
	// 	return
	// }

	url := "https://www.pro-football-reference.com/teams/" + team + "/"
	tableSelector := "#team_index"

	data, err := teamhistory.GetTeamHistory(url, tableSelector, year)

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
