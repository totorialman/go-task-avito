package models

import (
	"github.com/satori/uuid"
)

//easyjson:json
type Intake struct {
	ID     uuid.UUID `json:"id"`
	PVZID  uuid.UUID `json:"pvz_id"`
	Opened string    `json:"opened"`
	Closed *string   `json:"closed"`
}

//easyjson:json
type StartIntakeRequest struct {
	PVZID uuid.UUID `json:"pvz_id"`
}

//easyjson:json
type CloseIntakeRequest struct {
	PVZID uuid.UUID `json:"pvz_id"`
}
