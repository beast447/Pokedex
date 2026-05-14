package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"encoding/json"

	"github.com/beast447/pokedexcli/internal"
)

type commands struct {
	name        string
	description string
	callback    func(config *internal.Config, first string, second string) error
}

type pokemonStats struct {
	name   string
	hp     int
	attack int
	dead   bool
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
			name:        "pokedex",
			description: "Shows your current Pokedex",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows a captured pokemons stats like health, abilities , etc.",
			callback:    commandInspect,
		},
		"battle": {
			name:        "battle",
			description: "Fight two pokemon that you have stored in your pokedex",
			callback:    commandBattle,
		},
		"release": {
			name: "release",
			description: "Release a pokemon from your pokedex",
			callback: commandRelease,
		},
	}
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercase)
	slice := strings.Split(trimmed, " ")
	return slice
}

func getPokemonStats(config *internal.Config, firstPokemon string, secondPokemon string) ([]pokemonStats, error) {

	pokemonOne, exists := config.Pokedex[firstPokemon]
	if !exists {
		fmt.Printf("%v is not in your pokedex!\n", firstPokemon)
		return []pokemonStats{}, fmt.Errorf("%v is not in your pokedex!\n", firstPokemon)
	}
	fmt.Printf("%v is ready to fight!\n", pokemonOne.Name)

	pokemonTwo, exists := config.Pokedex[secondPokemon]
	if !exists {
		return []pokemonStats{}, fmt.Errorf("%v is not in your pokedex!\n", secondPokemon)
	}
	fmt.Printf("%v is ready to fight!\n", pokemonTwo.Name)

	var pokemonOneHp int
	for _, p := range pokemonOne.Stats {
		if p.Stat.Name == "hp" {
			pokemonOneHp = p.BaseStat
		}
	}
	var pokemonOneAttack int
	for _, p := range pokemonOne.Stats {
		if p.Stat.Name == "attack" {
			pokemonOneAttack = p.BaseStat
		}
	}
	var pokemonTwoHp int
	for _, p := range pokemonTwo.Stats {
		if p.Stat.Name == "hp" {
			pokemonTwoHp = p.BaseStat
		}
	}
	var pokemonTwoAttack int
	for _, p := range pokemonTwo.Stats {
		if p.Stat.Name == "attack" {
			pokemonTwoAttack = p.BaseStat
		}
	}

	sliceOfPokemonStats := []pokemonStats{
		{
			name:   firstPokemon,
			hp:     pokemonOneHp,
			attack: pokemonOneAttack,
			dead:   false,
		},
		{
			name:   secondPokemon,
			hp:     pokemonTwoHp,
			attack: pokemonTwoAttack,
			dead:   false,
		},
	}
	return sliceOfPokemonStats, nil
}

func commandExit(data *internal.Config, explore string, second string) error {
	fmt.Println("Saving your Pokedex...")
	save, err := json.Marshal(data.Pokedex)
	if err != nil{
		return err
	}
	if err := os.WriteFile(data.SavePath, (save), 0777); err != nil{
		return err
	}
	fmt.Print("Closing the Pokedex... Goodbye!")

	os.Exit(0)
	return nil
}

func commandHelp(data *internal.Config, explore string, second string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *internal.Config, explore string, second string) error {
	var err error
	savedPokedex := config.Pokedex
	if config.Next == "" {
		*config, err = internal.MakeInitialCall()
		config.Pokedex = savedPokedex
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

func commandMapBack(config *internal.Config, explore string, second string) error {
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

func commandExplore(config *internal.Config, explore string, second string) error {

	pokemon, err := internal.GetPokemonAtLocation(explore)
	if err != nil {
		return err
	}
	for _, i := range pokemon {
		fmt.Println(i.Name)
	}

	return nil
}

func commandCatch(config *internal.Config, selectedPokemon string, second string) error {
	userPokemon, err := internal.FindPokemon(selectedPokemon)
	if err != nil {
		return err
	}
	_, exists := config.Pokedex[selectedPokemon]
	if exists{
		fmt.Printf("You already caught %v!\n", selectedPokemon)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", selectedPokemon)
	userChance := rand.Intn(640)
	if userChance >= userPokemon.BaseExperience {
		config.Pokedex[selectedPokemon] = userPokemon
		fmt.Printf("%v was caught!\n", selectedPokemon)
	} else {
		fmt.Printf("%v escaped!\n", selectedPokemon)
	}
	return nil
}

func commandPokedex(config *internal.Config, selectedPokemon string, second string) error {

	if len(config.Pokedex) < 1 {
		fmt.Println("You dont have any pokemon! Get Catching!")
		return nil
	}

	fmt.Println("Your Pokedex:")

	for _, p := range config.Pokedex {
		fmt.Printf("- %v\n", p.Name)
	}

	return nil
}

func commandInspect(config *internal.Config, selectedPokemon string, second string) error {
	if len(config.Pokedex) < 1 {
		fmt.Println("You must catch a pokemon first to inspect it!")
		return nil
	}
	userPokemon, exists := config.Pokedex[selectedPokemon]
	if !exists {
		fmt.Printf("%s doesnt exist in your Pokedex yet\n", selectedPokemon)
		return nil
	} else {
		upperName := strings.ToUpper(userPokemon.Name)
		fmt.Printf("\n\n    %v\n\n", upperName)
		fmt.Println("------STATS--------")
		for _, stat := range userPokemon.Stats {
			fmt.Printf("- %v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("-----ABILITIES------")
		for _, a := range userPokemon.Abilities {
			fmt.Printf("- %v\n", a.Ability.Name)
		}
		fmt.Println("-------END---------\n\n")
	}
	return nil
}

func commandBattle(config *internal.Config, firstPokemon string, secondPokemon string) error {
	pokemonStats, err := getPokemonStats(config, firstPokemon, secondPokemon)
	if err != nil {
		fmt.Printf("Error fetching pokemon stats: %v\n", err)
	}
	pokemonOne := pokemonStats[0]
	pokemonTwo := pokemonStats[1]

	fmt.Printf("%v attack: %v\n", pokemonOne.name, pokemonOne.attack)
	fmt.Printf("%v HP: %v\n", pokemonTwo.name, pokemonTwo.hp)

	for !pokemonOne.dead || !pokemonTwo.dead {

		pokemonTwo.hp -= pokemonOne.attack
		if pokemonTwo.hp < 1 {
			pokemonTwo.dead = true
			pokemonTwo.hp = 0
			fmt.Printf("%v Wins!", pokemonOne.name)
			return nil
		}

		fmt.Printf("%v attacks %v with %v points of damage\n", pokemonOne.name, pokemonTwo.name, pokemonOne.attack)
		fmt.Printf("%v has %v points of health left\n", pokemonTwo.name, pokemonTwo.hp)

		pokemonOne.hp -= pokemonTwo.attack
		if pokemonOne.hp < 1 {
			pokemonOne.dead = true
			pokemonOne.hp = 0
			fmt.Printf("%v Wins!\n", pokemonTwo.name)
			return nil
		}

		fmt.Printf("%v attacks %v with %v points of damage\n", pokemonOne.name, pokemonTwo.name, pokemonOne.attack)
		fmt.Printf("%v has %v points of health left\n", pokemonTwo.name, pokemonTwo.hp)
	}
	return nil
}

func commandRelease (config *internal.Config, pokemonToRelease string, second string) error {
	_, exists := config.Pokedex[pokemonToRelease]
	if !exists {
		fmt.Printf("You dont have %v in your pokedex.\n", pokemonToRelease)
		return nil
	}
	delete(config.Pokedex, pokemonToRelease)
	fmt.Printf("%v was released!\n", pokemonToRelease)

	return nil
}
