package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/beast447/pokedexcli/internal"
)


func main() {

	config := &internal.Config{}
	scanner := bufio.NewScanner(os.Stdin)

	for{
	fmt.Print("Pokedex >")
	scanner.Scan()
	text := scanner.Text()
	cleanText := cleanInput(text)	
	command, exists := supportedCommands[cleanText[0]]
	if exists {
		err := command.callback(config)
		if err != nil{
			fmt.Printf("Error in callback function: %v", err)
			}
		}else {
			fmt.Print("Unkown command\n")
		}
}
}
