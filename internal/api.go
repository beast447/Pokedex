package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Pokedex            map[string]Pokemon
	SavePath           string
	CurrentAreaPokemon map[string]bool
}

type Location struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type LocPokemon struct {
	Name string
	Url  string
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

type Move struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
	Type  struct {
		Name string `json:"name"`
	} `json:"type"`
}

var cachedResult *Cache

func MakeInitialCall() (Config, error) {
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		return Config{}, fmt.Errorf("error on initial fetch: %v", err)
	}
	defer res.Body.Close()

	cacheData, err := io.ReadAll(res.Body)
	if err != nil {
		return Config{}, err
	}

	cachedResult = NewCache(5 * time.Minute)
	cachedResult.Add("https://pokeapi.co/api/v2/location-area/", cacheData)

	data := Config{}
	errr := json.Unmarshal(cacheData, &data)
	if errr != nil {
		return Config{}, fmt.Errorf("error when decoding data: %v", errr)
	}
	return data, nil
}

func GetNextLocation(config Config) (Config, error) {

	result, err := makeCallWithConfig(config, "Next")
	if err != nil {
		return Config{}, err
	}
	result.Pokedex = config.Pokedex
	return result, nil
}

func GetLastLocation(config Config) (Config, error) {

	result, err := makeCallWithConfig(config, "Previous")
	if err != nil {
		return Config{}, err
	}
	result.Pokedex = config.Pokedex
	return result, nil
}

func GetPokemonAtLocation(location string) ([]LocPokemon, error) {

	loc, err := makeCallWithString[Location](location, "Location")
	if err != nil {
		return []LocPokemon{}, err
	}

	result := make([]LocPokemon, 0, len(loc.PokemonEncounters))
	for _, p := range loc.PokemonEncounters {
		result = append(result, LocPokemon{Name: p.Pokemon.Name, Url: p.Pokemon.URL})
	}
	return result, nil

}

func FindPokemon(selectedPokemon string) (Pokemon, error) {
	pok, err := makeCallWithString[Pokemon](selectedPokemon, "Pokemon")
	if err != nil {
		return Pokemon{}, err
	}
	return pok, nil
}

func GetMove(name string) (Move, error) {
	if cachedResult == nil {
		cachedResult = NewCache(5 * time.Minute)
	}
	return makeCallWithString[Move](name, "Move")
}

func makeCallWithConfig(config Config, path string) (Config, error) {

	var url string
	switch path {
	case "Previous":
		url = config.Previous
	case "Next":
		url = config.Next
	}

	if cachedResult == nil {
		return Config{}, fmt.Errorf("cache not initialized, call map first")
	}
	stuff, exists := cachedResult.Get(url)
	if exists {
		data := Config{}
		err := json.Unmarshal(stuff, &data)
		if err != nil {
			return Config{}, err
		}
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return Config{}, fmt.Errorf("error fetching %s url: %v", path, err)
	}
	defer res.Body.Close()

	cacheData, err := io.ReadAll(res.Body)
	if err != nil {
		return Config{}, err
	}

	cachedResult.Add(url, cacheData)

	data := Config{}
	errr := json.Unmarshal(cacheData, &data)
	if errr != nil {
		return Config{}, fmt.Errorf("error when decoding data: %v", errr)
	}
	return data, nil
}

func makeCallWithString[T any](param string, path string) (T, error) {

	var result T
	var url string
	switch path {
	case "Location":
		url = "https://pokeapi.co/api/v2/location-area/" + param
	case "Pokemon":
		url = "https://pokeapi.co/api/v2/pokemon/" + param
	case "Move":
		url = "https://pokeapi.co/api/v2/move/" + param
	}

	if cachedResult != nil {
		if cached, exists := cachedResult.Get(url); exists {
			if err := json.Unmarshal(cached, &result); err != nil {
				return result, err
			}
			return result, nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	cacheData, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	if cachedResult != nil {
		cachedResult.Add(url, cacheData)
	}

	if err := json.Unmarshal(cacheData, &result); err != nil {
		return result, err
	}
	return result, nil
}
