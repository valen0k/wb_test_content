package weatherserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/valen0k/wb_test_content/internal/app/model"
	"github.com/valen0k/wb_test_content/internal/app/store"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
)

const APIURL = "https://api.openweathermap.org/data/2.5/forecast?appid=%s&lat=%f&lon=%f&mode=json&units=metric&lang=RU"

type WeatherServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *WeatherServer {
	return &WeatherServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (w *WeatherServer) Start(locationsFile string) error {
	if err := w.configureLogger(); err != nil {
		return err
	}

	locations, err := decodeLocationFile(locationsFile)
	if err != nil {
		return err
	}
	if err = w.configureStore(locations); err != nil {
		return err
	}
	defer w.store.Close()

	if err = w.uploadInfo(); err != nil {
		return err
	}

	w.configureRouter()

	w.logger.Infoln("starting weather server")

	return http.ListenAndServe(w.config.BindAddr, w.router)
}

func (w *WeatherServer) configureLogger() error {
	level, err := logrus.ParseLevel(w.config.LogLevel)
	if err != nil {
		return err
	}

	w.logger.SetLevel(level)

	return nil
}

func (w *WeatherServer) configureStore(locations map[string]model.Location) error {
	st := store.New(w.config.Store, locations)
	if err := st.Open(); err != nil {
		return err
	}

	w.store = st

	return nil
}

func (w *WeatherServer) uploadInfo() error {
	weathers, err := w.store.Weather().FindAll()
	if err != nil {
		return err
	}

	if len(weathers) != 0 {
		w.updateInfo(weathers)
	} else {
		w.downloadInfo()
	}

	return nil
}

func (w *WeatherServer) configureRouter() {
	w.router.HandleFunc("/cities", w.handleCities()).Methods("Get")
	w.router.HandleFunc("/cities/{city}", w.handleCity()).Methods("Get")
}

func (w *WeatherServer) handleCities() http.HandlerFunc {
	cities, err := w.store.Weather().FindAllCities()
	if err != nil {
		w.logger.Fatalln(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		w.respond(writer, request, http.StatusOK, cities)
	}
}

func (w *WeatherServer) handleCity() func(http.ResponseWriter, *http.Request) {
	allCities, err := w.store.Weather().FindAllCities()
	if err != nil {
		w.logger.Fatalln(err)
	}

	cities := make(map[string]struct{})
	for _, v := range allCities {
		if _, ok := cities[v]; !ok {
			cities[v] = struct{}{}
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		city := mux.Vars(request)["city"]
		if _, ok := cities[city]; !ok {
			w.error(writer, request, http.StatusBadRequest, errors.New("city not found"))
			return
		}

		weather, err := w.store.Weather().FindByCity(city)
		if err != nil {
			w.error(writer, request, http.StatusBadRequest, err)
			return
		}
		date := request.URL.Query().Get("date")
		if date == "" {

			res := model.RequestWeather{
				Country: weather.Data.City.Country,
				City:    city,
				Temp:    weather.Temp,
			}

			dates := make([]string, 0)
			for _, v := range weather.Data.List {
				dates = append(dates, v.DtTxt)
			}

			sort.Strings(dates)
			res.Dates = dates
			w.respond(writer, request, http.StatusOK, res)
			return
		}

		for i, list := range weather.Data.List {
			if list.DtTxt == date {
				weather.Data.List = weather.Data.List[i : i+1]
			}
		}

		w.respond(writer, request, http.StatusOK, weather)
		return
	}
}

func (w *WeatherServer) error(writer http.ResponseWriter, r *http.Request, code int, err error) {
	w.respond(writer, r, code, map[string]string{"error": err.Error()})
}

func (w *WeatherServer) respond(writer http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}

func (w *WeatherServer) downloadInfo() {
	wg := sync.WaitGroup{}
	for _, loc := range w.store.GetLocations() {
		wg.Add(1)
		go func(loc model.Location) {
			defer wg.Done()
			update, err := downloadInfoFromURL(loc, w.config.APIKey)
			if err != nil {
				w.logger.Error(err)
				return
			}
			err = w.store.Weather().Create(&update)
			if err != nil {
				w.logger.Error(err)
				return
			}
		}(loc)
	}

	wg.Wait()
}

func (w *WeatherServer) updateInfo(allWeather []model.Weather) {
	for _, v := range allWeather {
		go func(weather model.Weather) {
			if loc, ok := w.store.GetLocations()[weather.City]; ok {
				update, err := downloadInfoFromURL(loc, w.config.APIKey)
				if err != nil {
					w.logger.Error(err)
					return
				}
				err = w.store.Weather().Update(&update)
				if err != nil {
					w.logger.Error(err)
					return
				}
			}
		}(v)
	}
}

func decodeLocationFile(configFile string) (map[string]model.Location, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	locations := make([]model.Location, 0, 20)
	if err = json.Unmarshal(file, &locations); err != nil {
		return nil, err
	}

	res := make(map[string]model.Location)
	for _, v := range locations {
		if _, ok := res[v.City]; !ok {
			res[v.City] = v
		}
	}

	return res, nil
}

func downloadInfoFromURL(loc model.Location, APIKey string) (model.Weather, error) {
	var weather model.Weather

	resp, err := http.Get(fmt.Sprintf(
		APIURL,
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
