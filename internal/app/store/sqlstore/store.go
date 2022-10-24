package sqlstore

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/valen0k/wb_test_content/internal/app/model"
	"github.com/valen0k/wb_test_content/internal/app/store"
)

type Store struct {
	db         *sql.DB
	weatherRep *WeatherRepository
	locations  map[string]model.Location
}

func New(db *sql.DB, locations map[string]model.Location) *Store {
	return &Store{
		db:        db,
		locations: locations,
	}
}

func (s *Store) GetLocations() map[string]model.Location {
	return s.locations
}

func (s *Store) Weather() store.WeatherRepository {
	if s.weatherRep != nil {
		return s.weatherRep
	}

	s.weatherRep = &WeatherRepository{
		store: s,
	}

	return s.weatherRep
}
