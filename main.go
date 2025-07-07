package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexeychuk/pokedex_go/internal"
)

var baseApiUrl = "https://pokeapi.co/api/v2/"

type Config struct {
	Next     *string
	Previous *string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var nextUrl = baseApiUrl + "location-area"
	var config = Config{Next: &nextUrl, Previous: nil}

	cache := internal.NewCache(time.Duration(50 * time.Second))

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		userInputSlice := cleanInput(userInput)
		commandName := userInputSlice[0]

		if command, exists := commands[commandName]; !exists {
			fmt.Printf("Unknown command\n")
		} else {
			err := command.callback(&config, cache)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func cleanInput(text string) []string {

	return strings.Fields(strings.ToLower(text))
}
