package infrastructure

type PlayerView struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Country     string `json:"country"`
	Age         int    `json:"age"`
	Position    string `json:"position"`
	MarketValue int64  `json:"marketValue"`
}

type TeamView struct {
	ID             int64        `json:"id"`
	UserID         int64        `json:"userId"`
	Name           string       `json:"name"`
	Country        string       `json:"country"`
	Budget         int64        `json:"budget"`
	TotalTeamValue int64        `json:"totalTeamValue"`
	Players        []PlayerView `json:"players"`
}
