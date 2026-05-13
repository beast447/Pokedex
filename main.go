package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/beast447/pokedexcli/internal"
)



func main() {

	config := &internal.Config{Pokedex: make(map[string]internal.Pokemon)}
	scanner := bufio.NewScanner(os.Stdin)
	
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
		switch cleanText[0] {
		case "explore":
			if len(cleanText) < 2 {
				fmt.Println("Usage: explore <location>")
				continue
			}
			if err := command.callback(config, cleanText[1]); err != nil {
				fmt.Printf("Error in explore command: %v\n", err)
			}
		case "catch":
			if len(cleanText) < 2 {
				fmt.Println("Usage: catch <pokemon name>")
				continue
			}
			if err := command.callback(config, cleanText[1]); err != nil {
				fmt.Printf("Error in catch command: %v\n", err)
			}
		default:
			if err := command.callback(config, ""); err != nil {
				fmt.Printf("Error in callback function: %v\n", err)
			}
		}
	}
}
