package store

import (
	"github.com/valen0k/wb_test_content/internal/app/model"
)

type WeatherRepository struct {
	store *Store
}

func (r *WeatherRepository) Create(w *model.Weather) error {
	_, err := r.store.db.Exec("INSERT INTO weather (city, temp, data) VALUES ($1, $2, $3)", w.City, w.Temp, w.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *WeatherRepository) Update(w *model.Weather) error {
	_, err := r.store.db.Exec("UPDATE weather SET temp = $1, date = current_timestamp, data = $2 WHERE city = $3", w.Temp, w.Data, w.City)
	if err != nil {
		return err
	}

	return nil
}

func (r *WeatherRepository) FindAll() ([]model.Weather, error) {
	rows, err := r.store.db.Query("SELECT city, temp, date, data from weather order by city")
	if err != nil {
		return nil, err
	}

	res := make([]model.Weather, 0)
	for rows.Next() {
		var buf model.Weather
		rows.Scan(&buf.City, &buf.Temp, &buf.Date, &buf.Data)
		res = append(res, buf)
	}

	return res, nil
}

func (r *WeatherRepository) FindAllCities() ([]string, error) {
	rows, err := r.store.db.Query("SELECT city from weather order by city")
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, 20)
	for rows.Next() {
		var buf string
		rows.Scan(&buf)
		res = append(res, buf)
	}

	return res, nil
}

func (r *WeatherRepository) FindByCity(city string) (model.Weather, error) {
	var res model.Weather
	err := r.store.db.QueryRow(
		"SELECT city, temp, date, data from weather where city = $1",
		city,
	).Scan(&res.City, &res.Temp, &res.Date, &res.Data)
	if err != nil {
		return model.Weather{}, err
	}

	return res, nil
}
