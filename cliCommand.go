package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/Alexeychuk/pokedex_go/internal"
)

func commandExit(config *Config, cache *internal.Cache, _ []string, _ map[string]internal.PokemonResponse) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Error!!!")
}

func helpCallback(config *Config, cache *internal.Cache, _ []string, _ map[string]internal.PokemonResponse) error {
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

func mapCallback(config *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {

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

func mapbCallback(config *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {

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

func exploreLocationCallback(_ *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {
	if len(parameters) == 0 {
		fmt.Print("No location specified\n")
		return nil
	}
	name := parameters[0]
	mapData, err := internal.GetPokemonApiExploreLocation(name, cache)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", name)
	fmt.Print("Found Pokemon:\n")

	for _, pokemon := range mapData.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func catchCallback(_ *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {
	if len(parameters) == 0 {
		fmt.Print("No location specified\n")
		return nil
	}
	name := parameters[0]
	pokemonData, err := internal.GetPokemon(name, cache)

	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	if err != nil {
		return err
	}

	catchResult := rand.Intn(pokemonData.BaseExperience)

	if catchResult < pokemonData.BaseExperience-30 {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	pokedex[name] = pokemonData
	fmt.Printf("%s was caught!\n", name)
	return nil
}

func inspectCallback(_ *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {
	if len(parameters) == 0 {
		fmt.Print("No pokemon specified\n")
		return nil
	}

	name := parameters[0]
	pokemon, exists := pokedex[name]

	if !exists {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf(`Name: %s
Height: %d
Weight:  %d
`, pokemon.Name, pokemon.Height, pokemon.Weight)

	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("	-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, typeVal := range pokemon.Types {
		fmt.Printf("	-%s\n", typeVal.Type.Name)
	}

	return nil
}

func pokedexCallback(_ *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error {

	if len(pokedex) == 0 {
		fmt.Print("Youre pokedex is empty\n")
		return nil
	}

	fmt.Print("Your Pokedex:\n")
	for pokemon := range pokedex {
		fmt.Printf("- %s\n", pokemon)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *internal.Cache, parameters []string, pokedex map[string]internal.PokemonResponse) error
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
		"explore": {
			name:        "explore",
			description: "explores location, provides encountered pokemons",
			callback:    exploreLocationCallback,
		},
		"catch": {
			name:        "catch",
			description: "catch pokemon",
			callback:    catchCallback,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect pokemon",
			callback:    inspectCallback,
		},
		"pokedex": {
			name:        "pokedex",
			description: "show your pokedex",
			callback:    pokedexCallback,
		},
	}
}
