package weatherserver

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/valen0k/wb_test_content/internal/app/model"
	"github.com/valen0k/wb_test_content/internal/app/store"
	"net/http"
	"sort"
)

type server struct {
	logger *logrus.Logger
	router *mux.Router
	store  store.Store
}

func newServer(logger *logrus.Logger, store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
	}

	s.configureRouter()
	s.logger.Infoln("server create")

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/cities", s.handleCities()).Methods("Get")
	s.router.HandleFunc("/cities/{city}", s.handleCity()).Methods("Get")

	s.logger.Infoln("router create")
}

func (s *server) handleCities() http.HandlerFunc {
	cities, err := s.store.Weather().FindAllCities()
	if err != nil {
		s.logger.Fatalln(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		s.respond(writer, request, http.StatusOK, cities)
	}
}

func (s *server) handleCity() func(http.ResponseWriter, *http.Request) {
	allCities, err := s.store.Weather().FindAllCities()
	if err != nil {
		s.logger.Fatalln(err)
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
			err := errors.New("city not found")
			s.logger.Errorln(err)
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		weather, err := s.store.Weather().FindByCity(city)
		if err != nil {
			s.logger.Errorln(err)
			s.error(writer, request, http.StatusBadRequest, err)
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
			s.respond(writer, request, http.StatusOK, res)
			return
		}

		for i, list := range weather.Data.List {
			if list.DtTxt == date {
				weather.Data.List = weather.Data.List[i : i+1]
			}
		}

		s.respond(writer, request, http.StatusOK, weather)
		return
	}
}

func (s *server) error(writer http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(writer, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(writer http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
