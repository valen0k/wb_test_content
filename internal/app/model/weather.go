package model

import "time"

type Weather struct {
	City string          `json:"city"`
	Temp float64         `json:"temp"`
	Date time.Time       `json:"date"`
	Data ResponseWeather `json:"data"`
}
