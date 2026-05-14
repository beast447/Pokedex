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
		default:
			if err := command.callback(config, "", ""); err != nil {
				fmt.Printf("Error in callback function: %v\n", err)
			}
		}
	}
}
