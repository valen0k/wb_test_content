package weatherserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valen0k/wb_test_content/internal/app/model"
	"github.com/valen0k/wb_test_content/internal/app/store"
	"github.com/valen0k/wb_test_content/internal/app/store/sqlstore"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
)

const apiURL = "https://api.openweathermap.org/data/2.5/forecast?appid=%s&lat=%f&lon=%f&mode=json&units=metric&lang=RU"

func Start(config *Config, locFile string) error {
	locations, err := model.ReadLocFile(locFile)
	if err != nil {
		return err
	}

	logger, err := newLogger(config.LogLevel)
	if err != nil {
		return err
	}

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	logger.Infoln("connect database")

	st := sqlstore.New(db, locations)
	err = uploadWeather(st.Weather(), locations, config.APIKey)
	if err != nil {
		return err
	}

	srv := newServer(logger, st)
	logger.Infoln("server create")

	logger.Infoln("weather server start")
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newLogger(logLevel string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(level)

	return logger, nil
}

func uploadWeather(rep store.WeatherRepository, locations map[string]model.Location, apiKey string) error {
	weathers, err := rep.FindAll()
	if err != nil {
		return err
	}

	var funcRep func(*model.Weather) error
	if len(weathers) != 0 {
		funcRep = rep.Update
	} else {
		funcRep = rep.Create
	}

	downloadInfo(funcRep, locations, apiKey)
	return nil
}

func downloadInfo(f func(*model.Weather) error, locations map[string]model.Location, apiKey string) {
	wg := sync.WaitGroup{}
	for _, loc := range locations {
		wg.Add(1)
		go func(loc model.Location) {
			defer wg.Done()
			update, err := downloadInfoFromURL(loc, apiKey)
			if err != nil {
				log.Println(err)
				return
			}
			err = f(&update)
			if err != nil {
				log.Println(err)
				return
			}
		}(loc)
	}

	wg.Wait()
}

func downloadInfoFromURL(loc model.Location, APIKey string) (model.Weather, error) {
	var weather model.Weather

	resp, err := http.Get(fmt.Sprintf(
		apiURL,
		APIKey,
		loc.Lat,
		loc.Lon,
	))
	if err != nil {
		return model.Weather{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Weather{}, err
	}

	if err = json.Unmarshal(body, &weather.Data); err != nil {
		return model.Weather{}, err
	}

	var temp float64
	i := 0
	for _, v := range weather.Data.List {
		temp += v.Main.Temp
		i++
	}

	weather.City = loc.City
	weather.Temp = math.Round(temp/float64(i)*100) / 100

	return weather, nil
}
