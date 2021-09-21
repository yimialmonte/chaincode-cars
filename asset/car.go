package asset

// CarAsset ...
type Car struct {
	ID             string `json:"id"`
	Brand          string `json:"brand"`
	Owner          string `json:"owner"`
	TransfersCount int    `json:"transfersCount"`
}
