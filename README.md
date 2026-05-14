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
| `explore <location>` | List all Pokemon at a location and optionally catch one |
| `catch <pokemon>` | Attempt to catch a Pokemon |
| `pokedex` | View all your caught Pokemon |
| `inspect <pokemon>` | View a caught Pokemon's stats and abilities |
| `release <pokemon>` | Release a caught Pokemon from your Pokedex |
| `battle <pokemon1> <pokemon2>` | Battle two of your caught Pokemon |
| `compare <pokemon1> <pokemon2>` | Compare two of your caught Pokemon side by side |
| `exit` | Save and quit |

## Interactive UI

The CLI uses [pterm](https://github.com/pterm/pterm) for a rich terminal experience:

- **`explore`** — opens an interactive menu to pick a Pokemon to catch, with an option to stay in the area and try again
- **`pokedex`** — opens an interactive list; select a Pokemon to inspect its full stats
- **`catch`** — shows a spinner while throwing the Pokeball and renders the Pokemon's sprite as ASCII art on a successful catch
- **`inspect`** — displays stats and abilities in formatted tables with the Pokemon's sprite

## Catching

Each Pokemon has a `base_experience` value from the PokeAPI. When you throw a Pokeball, a random number between 0 and 639 is generated — if it's greater than or equal to the Pokemon's base experience, the catch succeeds. Rarer, stronger Pokemon are harder to catch.

## Battling

Battles are turn-based. Each Pokemon attacks using one of its actual moves (falling back to its `attack` stat if no damaging move is found). HP bars are rendered each round with color-coded health levels. The first Pokemon to reach 0 HP loses.

## Comparing

`compare <pokemon1> <pokemon2>` displays both Pokemon's stats and abilities side by side in a panel layout, making it easy to size up your team.

## Project Structure

```
.
├── main.go          # REPL loop, startup banner, and command dispatch
├── repl.go          # Command definitions, handlers, and rendering
├── repl_test.go     # Tests for input cleaning and helper logic
└── internal/
    ├── api.go       # PokeAPI client and data types
    └── pokecache.go # Thread-safe in-memory cache
```
