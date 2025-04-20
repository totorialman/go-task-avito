package models

import (
	"github.com/satori/uuid"
)

//easyjson:json
type Product struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Price    int       `json:"price"`
	IntakeID uuid.UUID `json:"intake_id"`
	Created  string    `json:"created"`
}

//easyjson:json
type AddProductRequest struct {
	Name  string    `json:"name"`
	Price int       `json:"price"`
	PVZID uuid.UUID `json:"pvz_id"`
}

//easyjson:json
type DeleteLastProductRequest struct {
	PVZID uuid.UUID `json:"pvz_id"`
}
