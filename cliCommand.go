package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Alexeychuk/pokedex_go/internal"
)

func commandExit(config *Config, cache *internal.Cache) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Error!!!")
}

func helpCallback(config *Config, cache *internal.Cache) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}

	return nil
}

func getCachedOrNewMapData(cache *internal.Cache, url string) (internal.MapResponse, error) {
	cachedData, isCached := cache.Get(url)
	var mapData internal.MapResponse

	if isCached {
		reader := bytes.NewReader(cachedData)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&mapData); err != nil {
			return internal.MapResponse{}, err
		}
	} else {
		data, err := internal.GetPokemonApiLocations(url, cache)
		if err != nil {
			return internal.MapResponse{}, err
		}
		mapData = data
	}
	return mapData, nil
}

func mapCallback(config *Config, cache *internal.Cache) error {

	if config.Next == nil {
		fmt.Print("You are on last page\n")
		return nil
	}

	mapData, err := getCachedOrNewMapData(cache, *config.Next)
	if err != nil {
		return err
	}

	config.Previous = config.Next
	next := &mapData.Next
	config.Next = next

	for _, r := range mapData.Results {
		fmt.Printf("%s\n", r.Name)
	}

	return nil
}

func mapbCallback(config *Config, cache *internal.Cache) error {

	if config.Previous == nil {
		fmt.Print("You are on first page\n")
		return nil
	}

	mapData, err := getCachedOrNewMapData(cache, *config.Previous)
	if err != nil {
		return err
	}

	config.Next = config.Previous
	config.Previous = mapData.Previous

	for _, r := range mapData.Results {
		fmt.Printf("%s\n", r.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *internal.Cache) error
}

// Declare the map variable
var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCallback,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "gives 20 locations",
			callback:    mapCallback,
		},
		"mapb": {
			name:        "mapb",
			description: "gives 20 prev locations",
			callback:    mapbCallback,
		},
	}
}
