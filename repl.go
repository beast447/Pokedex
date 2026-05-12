package main

import (
	"strings"
	"fmt"
	"os"

	"github.com/beast447/pokedexcli/internal"
)

type commands struct{
	name string
	description string
	callback func(config *internal.Config)error
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
			name: "map",
			description: "Displays a list of 20 locations in the pokemon world",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Gets the previous list of 20 locations",
			callback: commandMapBack,
		},
	}
}

func cleanInput(text string) []string{
	lowercase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercase)
	slice := strings.Split(trimmed, " ")
	return slice
}

func commandExit(data *internal.Config) error{
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(data *internal.Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands{
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *internal.Config) error {
	var err error
	if  config.Next == "" {
		*config, err = internal.MakeInitialCall()
		
	} else{
		*config, err = internal.GetNextLocation(*config)
	}
	if err != nil{
		return err
	}
	for _, location := range config.Results{
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapBack (config *internal.Config) error {
	var err error
	if config.Previous == "" {
	 fmt.Println("you're on the first page")
	return nil
	} else{
	*config, err = internal.GetLastLocation(*config)
	}
	if err != nil{
	return err
	}
	for _, location := range config.Results{
		fmt.Println(location.Name)
	}
	return nil
}




