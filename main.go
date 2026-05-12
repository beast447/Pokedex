package main

import (
	"bufio"
	"fmt"
	"os"

)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for{
	fmt.Print("Pokedex >")
	scanner.Scan()
	text := scanner.Text()
	cleanText := cleanInput(text)	
	fmt.Printf("Your command was: %v\n", cleanText[0])
	}


}
