package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beast447/pokedexcli/internal"
	"github.com/pterm/pterm"
)



func main() {
	
	config := &internal.Config{
		Pokedex:            make(map[string]internal.Pokemon),
		CurrentAreaPokemon: make(map[string]bool),
	}
	scanner := bufio.NewScanner(os.Stdin)
	
	cwd, err := os.Getwd()
	if err != nil{
		log.Fatal(err)
	}
	savePath := filepath.Join(cwd, "save.json")
	config.SavePath = savePath
	
	if _, err := os.Stat(savePath); err != nil{
		fmt.Printf("No save file detected")
	} else{
		saveData, err := os.ReadFile(savePath)
		if err != nil{
			log.Fatal(err)
		}
		if err := json.Unmarshal(saveData, &config.Pokedex); err != nil{
			log.Fatal(err)	
		}
	}

	
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleanText := cleanInput(text)
		if cleanText[0] == "" {
			continue
		}
		command, exists := supportedCommands[cleanText[0]]
		if !exists {
			fmt.Print("Unknown command\n")
			continue
		}
		fmt.Print("\033[H\033[2J")
		pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("P O K E D E X"))
		fmt.Println()
		switch cleanText[0] {
		case "explore":
			if len(cleanText) < 2 {
				fmt.Println("Usage: explore <location>")
				continue
			}
			if err := command.callback(config, cleanText[1], ""); err != nil {
				fmt.Printf("Error in explore command: %v\n", err)
			}
		case "catch":
			if len(cleanText) < 2 {
				fmt.Println("Usage: catch <pokemon name>")
				continue
			}
			if err := command.callback(config, cleanText[1], ""); err != nil {
				fmt.Printf("Error in catch command: %v\n", err)
			}
		case "inspect":
			if len(cleanText) < 2{
				fmt.Println("Usage: inspect <Pokemon Name> Must be in your pokedex to inspect!")
				continue
			}
			if err := command.callback(config, cleanText[1], ""); err != nil{
				fmt.Printf("Error in inspect command: %v\n", err)
			}
		case "battle":
			if len(cleanText) < 3{
				fmt.Println("Usage: battle <first pokemon> <second pokemon> Both must be in your pokedex!")
			}
			if err := command.callback(config, cleanText[1], cleanText[2]); err != nil{
				fmt.Printf("Error in battle callback: %v", err)
			}
		case "release":
			if len(cleanText) < 2 {
				fmt.Println("Usage: release <pokemon name>")
				continue
			}
			if err := command.callback(config, cleanText[1], ""); err != nil {
				fmt.Printf("Error in release command: %v\n", err)
			}
		case "compare":
			first, second := "", ""
			if len(cleanText) >= 2 {
				first = cleanText[1]
			}
			if len(cleanText) >= 3 {
				second = cleanText[2]
			}
			if err := command.callback(config, first, second); err != nil {
				fmt.Printf("Error in compare command: %v\n", err)
			}
		default:
			if err := command.callback(config, "", ""); err != nil {
				fmt.Printf("Error in callback function: %v\n", err)
			}
		}
	}
}
