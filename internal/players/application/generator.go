package application

import (
	"math/rand"
	"time"

	"github.com/sergiojaa/soccer-manager-api/internal/players/domain"
)

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

func GeneratePlayers(teamID int64) []domain.PlayerInput {
	rand.Seed(time.Now().UnixNano())

	var players []domain.PlayerInput

	addPlayers := func(count int, position string) {
		for i := 0; i < count; i++ {
			player := domain.PlayerInput{
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
