package models

//easyjson:json
type PVZFilterRequest struct {
	StartDate string `json:"start_date"` // ISO8601
	EndDate   string `json:"end_date"`   // ISO8601
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

//easyjson:json
type PVZWithIntakes struct {
	PVZ     PVZ      `json:"pvz"`
	Intakes []Intake `json:"intakes"`
}

//easyjson:json
type PVZListResponse struct {
	Items []PVZWithIntakes `json:"items"`
}
