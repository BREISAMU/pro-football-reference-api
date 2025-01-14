# Pro Football Reference API [www.pro-football-reference.com/]
Go api servicing GET requests for Pro Football Reference tables.
<br/>
<br/>
<b>DNS Address </b>
```
BASE_URL=ec2-18-118-33-121.us-east-2.compute.amazonaws.com
```
<b>Functioning Request Routes</b> ( see [handlers](https://github.com/BREISAMU/pro-football-reference-api/tree/main/handlers) for further documentation ):
<br/>
<h1>/team</h1>
<ul>
    <li> /?team=TEAM_NAME&year=YEAR</li>
    <li> /draft/?team=TEAM_NAME&year=YEAR</li>
    <li> /offensiveStats?team=TEAM_NAME&year=YEAR</li>
    <li> /defensiveStats?team=TEAM_NAME&year=YEAR</li>
    <li> /offensiveRankings?team=TEAM_NAME&year=YEAR</li>
    <li> /defensiveRankings?team=TEAM_NAME&year=YEAR</li>
</ul>

## Example usage   
![plot](./images/rawTable.png)
```
curl "ec2-18-118-33-121.us-east-2.compute.amazonaws.com/team/?team=atl/year=2015"
```
<br/> -->
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

# Disclaimer / Credit
All data used in this project is sourced from [Sports Reference LLC](https://www.sports-reference.com/?utm_source=sr&utm_medium=sr_xsite&utm_campaign=2023_01_srnav). Use of this API must comply with all conditions set forth by their [terms and conditions](https://www.sports-reference.com/termsofuse.html?__hstc=223721476.0095f109d09965649dce377d8f2cfafe.1734971558705.1736797649244.1736804227068.13&__hssc=223721476.35.1736804227068&__hsfp=424417210) and [scraping/bot guidelines](https://www.sports-reference.com/bot-traffic.html). PLEASE do not use this service to produce archives of SRL data or any other action prohibited by the terms and conditions of SRL. All use of this API stems from one instance, and is collectively subject to the individual rate limit put forth by SRL. Take reasonable measures to stay below the rate limit of 20 requests per minute or this service may be taken down.
