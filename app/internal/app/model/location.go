package model

import (
	"encoding/json"
	"os"
)

type Location struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

func ReadLocFile(locFile string) (map[string]Location, error) {
	file, err := os.ReadFile(locFile)
	if err != nil {
		return nil, err
	}

	locations := make([]Location, 0, 20)
	if err = json.Unmarshal(file, &locations); err != nil {
		return nil, err
	}

	res := make(map[string]Location)
	for _, v := range locations {
		if _, ok := res[v.City]; !ok {
			res[v.City] = v
		}
	}

	return res, nil
}
