package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexeychuk/pokedex_go/internal"
)

type Config struct {
	Next     *string
	Previous *string
}

type Pokedex map[string]internal.PokemonResponse

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var config = Config{Next: &internal.NextUrl, Previous: nil}

	cache := internal.NewCache(time.Duration(5 * time.Second))
	pokedex := Pokedex{}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		userInputSlice := cleanInput(userInput)
		if len(userInputSlice) == 0 {
			fmt.Printf("Unknown command\n")
			continue
		}

		commandName := userInputSlice[0]
		parameters := userInputSlice[1:]

		if command, exists := commands[commandName]; !exists {
			fmt.Printf("Unknown command\n")
		} else {
			err := command.callback(&config, cache, parameters, pokedex)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func cleanInput(text string) []string {

	return strings.Fields(strings.ToLower(text))
}
