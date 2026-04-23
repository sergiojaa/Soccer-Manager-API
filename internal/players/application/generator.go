package application

import (
	"math/rand"
	"time"
)

type PlayerInput struct {
	TeamID    int64
	FirstName string
	LastName  string
	Country   string
	Age       int
	Position  string
}

const (
	GKCount  = 3
	DEFCount = 6
	MIDCount = 6
	ATTCount = 5
)

const (
	PositionGK  = "GK"
	PositionDEF = "DEF"
	PositionMID = "MID"
	PositionATT = "ATT"
)

var firstNames = []string{"John", "Alex", "Mike", "Leo", "David"}
var lastNames = []string{"Smith", "Brown", "Taylor", "Wilson", "Johnson"}
var countries = []string{"England", "Spain", "Germany", "France", "Italy"}

func GeneratePlayers(teamID int64) []PlayerInput {
	rand.Seed(time.Now().UnixNano())

	var players []PlayerInput

	addPlayers := func(count int, position string) {
		for i := 0; i < count; i++ {
			player := PlayerInput{
				TeamID:    teamID,
				FirstName: randomFrom(firstNames),
				LastName:  randomFrom(lastNames),
				Country:   randomFrom(countries),
				Age:       rand.Intn(23) + 18,
				Position:  position,
			}
			players = append(players, player)
		}
	}

	addPlayers(GKCount, PositionGK)
	addPlayers(DEFCount, PositionDEF)
	addPlayers(MIDCount, PositionMID)
	addPlayers(ATTCount, PositionATT)

	return players
}

func randomFrom(list []string) string {
	return list[rand.Intn(len(list))]
}
