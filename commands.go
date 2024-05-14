package main

import (
	"fmt"
	"os"

	"github.com/Ekwinder/pokedexcli/internal/pokeapi"
)

var mapOrder = []string{"help", "exit", "map", "mapb"}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		mapOrder[0]: {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		mapOrder[1]: {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		mapOrder[2]: {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations, and so on",
			callback:    commandMap,
		},
		mapOrder[3]: {
			name:        "mapb",
			description: `Similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations. Returns Error if already on the first page`,
			callback:    commandMapB,
		},
	}
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for _, k := range mapOrder {
		fmt.Printf("%s: %s\n", getCommands()[k].name, getCommands()[k].description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMap() error {
	pokeapi.GetMap(false)
	return nil
}

func commandMapB() error {
	pokeapi.GetMap(true)
	return nil
}
