package main

import (
	"strings"
)

func cleanInput(text string) []string{
	lowercase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercase)
	slice := strings.Split(trimmed, " ")
	return slice
}


