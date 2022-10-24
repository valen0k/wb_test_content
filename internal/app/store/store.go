package store

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/valen0k/wb_test_content/internal/app/model"
)

type Store struct {
	config     *Config
	db         *sql.DB
	weatherRep *WeatherRepository
	locations  map[string]model.Location
}

func New(config *Config, locations map[string]model.Location) *Store {
	return &Store{
		config:    config,
		locations: locations,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("pgx", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) GetLocations() map[string]model.Location {
	return s.locations
}

func (s *Store) Weather() *WeatherRepository {
	if s.weatherRep != nil {
		return s.weatherRep
	}

	s.weatherRep = &WeatherRepository{
		store: s,
	}

	return s.weatherRep
}
