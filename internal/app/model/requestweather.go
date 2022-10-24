package model

type RequestWeather struct {
	Country string   `json:"country"`
	City    string   `json:"city"`
	Temp    float64  `json:"temp"`
	Dates   []string `json:"dates"`
}
