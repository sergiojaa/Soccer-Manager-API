package infrastructure

type MarketListingView struct {
	ID           int64  `json:"id"`
	PlayerID     int64  `json:"playerId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Country      string `json:"country"`
	Age          int    `json:"age"`
	Position     string `json:"position"`
	MarketValue  int64  `json:"marketValue"`
	AskingPrice  int64  `json:"askingPrice"`
	SellerTeamID int64  `json:"sellerTeamId"`
	SellerTeam   string `json:"sellerTeam"`
}
