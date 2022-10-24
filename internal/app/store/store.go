package store

type Store interface {
	Weather() WeatherRepository
}
