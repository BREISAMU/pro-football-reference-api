# Pro Football Reference     API [https://www.pro-football-reference.com/]
Service GET requests for Pro Football Reference tables. Return in JSON format.

<b>Current functionality:</b>:
<br/>
- /teams/TEAM_NAME
<br/>

### Example usage   [ GET teamHistory?team=atl/year=2015 ]
![plot](./images/rawTable.png)
<br/><br/> -->
<br/><br/>
```
{
    "year": 2015,
    "league": "NFL",
    "team": "Atlanta Falcons",
    "wins": 8,
    "losses": 8,
    "ties": 0,
    "divisionFinish": 2,
    "playoffExitRound": 0,
    "pointsFor": 339,
    "pointsAgainst": 345,
    "pointsDif": -6,
    "headCoaches": "Quinn",
    "bestPlayerAv": "Jones",
    "bestPlayerPasser": "Ryan",
    "bestPlayerRusher": "Freeman",
    "bestPlayerReceiver": "Jones",
    "offRankPts": 21,
    "offRankYds": 7,
    "defRankPts": 14,
    "defRankYds": 16,
    "takeawayRank": 27,
    "pointsDifRank": 15,
    "yardsDifRank": 10,
    "teamsInLeague": 32,
    "marginOfVictory": -0.4,
    "strengthOfSchedule": -3.4,
    "srs": -3.8,
    "offensiveSrs": -4,
    "defensiveSrs": 0.3
}
```

#
(https://www.pro-football-reference.com/)
