# Pokedex CLI

A command-line Pokedex built in Go that lets you explore Pokemon locations, catch Pokemon, and manage your collection. Data is sourced from the [PokeAPI](https://pokeapi.co/).

## Features

- Browse paginated Pokemon location areas
- Explore specific locations to see which Pokemon appear there
- Attempt to catch Pokemon (success is based on their base experience)
- View your caught Pokemon in your personal Pokedex
- In-memory response caching with a 5-minute TTL to avoid redundant API calls

## Installation

```bash
git clone https://github.com/beast447/pokedexcli
cd pokedexcli
go build -o pokedexcli
./pokedexcli
```

## Commands

| Command | Description |
|---|---|
| `help` | Display all available commands |
| `map` | Show the next 20 Pokemon location areas |
| `mapb` | Go back to the previous 20 location areas |
| `explore <location>` | List all Pokemon found at a given location area |
| `catch <pokemon>` | Attempt to catch a Pokemon by name |
| `pokedex` | View all Pokemon you've caught |
| `exit` | Quit the program |

## Example Usage

```
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore pastoria-city-area
tentacool
tentacruel
shellos
...

Pokedex > catch shellos
Throwing a Pokeball at shellos...
shellos was caught!

Pokedex > pokedex
Your Pokedex:
- shellos
```

## How Catching Works

Each Pokemon has a `base_experience` value from the PokeAPI. When you throw a Pokeball, a random number between 0 and 639 is generated. If that number is greater than or equal to the Pokemon's base experience, the catch succeeds. Rarer, stronger Pokemon are harder to catch.

## Project Structure

```
.
├── main.go          # REPL loop and command dispatch
├── repl.go          # Command definitions and handlers
└── internal/
    ├── api.go       # PokeAPI client and data types
    └── pokecache.go # Thread-safe in-memory cache
```
