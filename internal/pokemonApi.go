package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MapResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetPokemonApiLocations(url string, cache *Cache) (MapResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return MapResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return MapResponse{}, fmt.Errorf("status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return MapResponse{}, err
	}

	var resData MapResponse

	if err = json.Unmarshal(data, &resData); err != nil {
		return MapResponse{}, err
	}

	cache.Add(url, data)

	return resData, nil
}
