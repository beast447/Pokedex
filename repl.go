package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/beast447/pokedexcli/internal"
)

type commands struct {
	name        string
	description string
	callback    func(config *internal.Config, explore string) error
}

var supportedCommands map[string]commands

func init() {
	supportedCommands = map[string]commands{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of 20 locations in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets the previous list of 20 locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all the pokemon at a selected location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon by name",
			callback:    commandCatch,
		},
		"pokedex": {
			name: "pokedex",
			description: "Shows your current Pokedex",
			callback: commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercase)
	slice := strings.Split(trimmed, " ")
	return slice
}

func commandExit(data *internal.Config, explore string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(data *internal.Config, explore string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *internal.Config, explore string) error {
	var err error
	if config.Next == "" {
		*config, err = internal.MakeInitialCall()

	} else {
		*config, err = internal.GetNextLocation(*config)
	}
	if err != nil {
		return err
	}
	for _, location := range config.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapBack(config *internal.Config, explore string) error {
	var err error
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		*config, err = internal.GetLastLocation(*config)
	}
	if err != nil {
		return err
	}
	for _, location := range config.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(config *internal.Config, explore string) error {

	pokemon, err := internal.GetPokemonAtLocation(explore)
	if err != nil {
		return err
	}
	for _, i := range pokemon {
		fmt.Println(i.Name)
	}

	return nil
}

func commandCatch(config *internal.Config, selectedPokemon string) error {
	userPokemon, err := internal.FindPokemon(selectedPokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", selectedPokemon)
	userChance := rand.Intn(640)
	if userChance >= userPokemon.BaseExperience {
		config.Pokedex[selectedPokemon] = userPokemon
		fmt.Printf("%v was caught!\n", selectedPokemon)
	} else {
		fmt.Printf("%v escaped!\n")
	}
	return nil
}

func commandPokedex (config *internal.Config, selectedPokemon string) error{
	
	if len(config.Pokedex) < 1 {
		fmt.Println("You dont have any pokemon! Get Catching!")
		return nil
	}

	fmt.Println("Your Pokedex:")

	for _, p := range config.Pokedex{
		fmt.Printf("- %v\n", p.Name)
	}

return nil
}
