package main

import (
	"fmt"
	"os"

	"github.com/Ekwinder/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(names ...string) error
}

var mapOrder = [5]string{"help", "exit", "map", "mapb", "explore"}

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
		mapOrder[4]: {
			name:        "explore",
			description: `See a list of all the Pokémon in a given area.`,
			callback:    explore,
		},
	}
}

func commandHelp(names ...string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for _, k := range mapOrder {
		fmt.Printf("%s: %s\n", getCommands()[k].name, getCommands()[k].description)
	}
	return nil
}

func commandExit(names ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(names ...string) error {
	pokeapi.GetMap(false)
	return nil
}

func commandMapB(names ...string) error {
	pokeapi.GetMap(true)
	return nil
}

func explore(names ...string) error {
	if len(names) > 0 {
		pokeapi.Explore(names[0])
		return nil
	} else {
		return fmt.Errorf("Please provide a valid name")
	}
}
