package store

import "github.com/valen0k/wb_test_content/internal/app/model"

type WeatherRepository interface {
	Create(*model.Weather) error
	Update(*model.Weather) error
	FindAll() ([]model.Weather, error)
	FindAllCities() ([]string, error)
	FindByCity(city string) (model.Weather, error)
}
