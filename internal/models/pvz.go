package models

import (
	"github.com/satori/uuid"
)

type City string

const (
	CityMoscow City = "Москва"
	CitySpb    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

//easyjson:json
type PVZ struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	City    City      `json:"city"`
	Address string    `json:"address"`
	Created string    `json:"created"`
}

//easyjson:json
type CreatePVZRequest struct {
	Name    string `json:"name"`
	City    City   `json:"city"`
	Address string `json:"address"`
}
