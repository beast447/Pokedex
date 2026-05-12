package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


func MakeInitialCall() (Config, error){
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil{
		return Config{}, fmt.Errorf("error on initial fetch: %v", err)
	}
	defer res.Body.Close()

	data := Config{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil{
		return Config{}, fmt.Errorf("error when decoding date: %v", err)
	}
	
	cache := NewCache(5 * time.Second)
	for _, item := range data.Results{
		cache.Add(item, )
	}

	return  data, nil
}

func GetNextLocation(config Config) (Config, error) {
	res, err := http.Get(config.Next)
	if err != nil{
		return Config{}, fmt.Errorf("error fetching next url: %v", err)
	}
	defer res.Body.Close()
	
	data := Config{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil{
		return Config{}, fmt.Errorf("error decoding next url: %v", err)
	}

	return data, nil
}

func GetLastLocation(config Config) (Config, error) {
	res, err := http.Get(config.Previous)
	if err != nil{
		return Config{}, fmt.Errorf("error fetching previous url: %v", err)
	}
	defer res.Body.Close()
	
	data := Config{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil{
		return Config{}, fmt.Errorf("error decoding previous url: %v", err)
	}

	return data, nil}
