# Pokedex CLI

A command-line Pokedex built in Go that lets you explore Pokemon locations, catch Pokemon, and manage your collection. Data is sourced from the [PokeAPI](https://pokeapi.co/).

## Features

- Browse paginated Pokemon location areas
- Explore specific locations to see which Pokemon appear there
- Attempt to catch Pokemon (success is based on their base experience)
- View your caught Pokemon in your personal Pokedex
- Inspect a caught Pokemon's stats and abilities
- Battle two of your caught Pokemon against each other
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
| `inspect <pokemon>` | View a caught Pokemon's stats and abilities |
| `battle <pokemon1> <pokemon2>` | Simulate a turn-based battle between two caught Pokemon |
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

Pokedex > inspect shellos
shellos

------STATS--------
- hp: 76
- attack: 48
...
-----ABILITIES------
- sticky-hold
- storm-drain
-------END---------

Pokedex > battle shellos pikachu
shellos is ready to fight!
pikachu is ready to fight!
shellos attacks pikachu with 48 points of damage
pikachu has 7 points of health left
...
shellos Wins!
```

## How Catching Works

Each Pokemon has a `base_experience` value from the PokeAPI. When you throw a Pokeball, a random number between 0 and 639 is generated. If that number is greater than or equal to the Pokemon's base experience, the catch succeeds. Rarer, stronger Pokemon are harder to catch.

## How Battling Works

Battles are turn-based. Each Pokemon attacks using its `attack` base stat, dealing that much damage to the opponent's `hp` each round. The first Pokemon to reach 0 HP loses. Both Pokemon must already be in your Pokedex to battle.

## Project Structure

```
.
├── main.go          # REPL loop and command dispatch
├── repl.go          # Command definitions and handlers
└── internal/
    ├── api.go       # PokeAPI client and data types
    └── pokecache.go # Thread-safe in-memory cache
```
