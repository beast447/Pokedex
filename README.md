# Pokedex CLI

A command-line Pokedex built in Go. Explore the Pokemon world, catch and manage your collection, and battle your Pokemon against each other. Data is sourced from the [PokeAPI](https://pokeapi.co/). Your Pokedex is saved automatically when you exit.

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
| `map` | Show the next 20 location areas |
| `mapb` | Go back to the previous 20 location areas |
| `explore <location>` | List all Pokemon found at a location |
| `catch <pokemon>` | Attempt to catch a Pokemon |
| `pokedex` | View all your caught Pokemon |
| `inspect <pokemon>` | View a caught Pokemon's stats and abilities |
| `release <pokemon>` | Release a caught Pokemon from your Pokedex |
| `battle <pokemon1> <pokemon2>` | Battle two of your caught Pokemon |
| `exit` | Save and quit |

## Catching

Each Pokemon has a `base_experience` value from the PokeAPI. When you throw a Pokeball, a random number between 0 and 639 is generated — if it's greater than or equal to the Pokemon's base experience, the catch succeeds. Rarer, stronger Pokemon are harder to catch.

## Battling

Battles are turn-based. Each Pokemon attacks using its `attack` stat, dealing that much damage to the opponent's `hp` each round. The first to reach 0 HP loses.

## Project Structure

```
.
├── main.go          # REPL loop and command dispatch
├── repl.go          # Command definitions and handlers
└── internal/
    ├── api.go       # PokeAPI client and data types
    └── pokecache.go # Thread-safe in-memory cache
```
