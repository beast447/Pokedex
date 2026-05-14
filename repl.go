package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beast447/pokedexcli/internal"
	"github.com/pterm/pterm"
	"github.com/qeesung/image2ascii/convert"
)

type commands struct {
	name        string
	description string
	callback    func(config *internal.Config, first string, second string) error
}

type pokemonStats struct {
	name   string
	hp     int
	attack int
	dead   bool
}

var supportedCommands map[string]commands

func init() {
	supportedCommands = map[string]commands{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of 20 locations in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets the previous list of 20 locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all the pokemon at a selected location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon by name",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows your current Pokedex",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows a captured pokemons stats like health, abilities , etc.",
			callback:    commandInspect,
		},
		"battle": {
			name:        "battle",
			description: "Fight two pokemon that you have stored in your pokedex",
			callback:    commandBattle,
		},
		"release": {
			name: "release",
			description: "Release a pokemon from your pokedex",
			callback: commandRelease,
		},
		"compare": {
			name:        "compare",
			description: "Compare two pokemon from your pokedex side by side",
			callback:    commandCompare,
		},
	}
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercase)
	slice := strings.Split(trimmed, " ")
	return slice
}

func getPokemonStats(config *internal.Config, firstPokemon string, secondPokemon string) ([]pokemonStats, error) {

	pokemonOne, exists := config.Pokedex[firstPokemon]
	if !exists {
		fmt.Printf("%v is not in your pokedex!\n", firstPokemon)
		return []pokemonStats{}, fmt.Errorf("%v is not in your pokedex!\n", firstPokemon)
	}
	fmt.Printf("%v is ready to fight!\n", pokemonOne.Name)

	pokemonTwo, exists := config.Pokedex[secondPokemon]
	if !exists {
		return []pokemonStats{}, fmt.Errorf("%v is not in your pokedex!\n", secondPokemon)
	}
	fmt.Printf("%v is ready to fight!\n", pokemonTwo.Name)

	var pokemonOneHp int
	for _, p := range pokemonOne.Stats {
		if p.Stat.Name == "hp" {
			pokemonOneHp = p.BaseStat
		}
	}
	var pokemonOneAttack int
	for _, p := range pokemonOne.Stats {
		if p.Stat.Name == "attack" {
			pokemonOneAttack = p.BaseStat
		}
	}
	var pokemonTwoHp int
	for _, p := range pokemonTwo.Stats {
		if p.Stat.Name == "hp" {
			pokemonTwoHp = p.BaseStat
		}
	}
	var pokemonTwoAttack int
	for _, p := range pokemonTwo.Stats {
		if p.Stat.Name == "attack" {
			pokemonTwoAttack = p.BaseStat
		}
	}

	sliceOfPokemonStats := []pokemonStats{
		{
			name:   firstPokemon,
			hp:     pokemonOneHp,
			attack: pokemonOneAttack,
			dead:   false,
		},
		{
			name:   secondPokemon,
			hp:     pokemonTwoHp,
			attack: pokemonTwoAttack,
			dead:   false,
		},
	}
	return sliceOfPokemonStats, nil
}

func commandExit(data *internal.Config, explore string, second string) error {
	fmt.Println("Saving your Pokedex...")
	save, err := json.Marshal(data.Pokedex)
	if err != nil{
		return err
	}
	if err := os.WriteFile(data.SavePath, (save), 0777); err != nil{
		return err
	}
	fmt.Print("Closing the Pokedex... Goodbye!")

	os.Exit(0)
	return nil
}

func commandHelp(data *internal.Config, explore string, second string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *internal.Config, explore string, second string) error {
	var err error
	savedPokedex := config.Pokedex
	savedPath := config.SavePath
	savedArea := config.CurrentAreaPokemon
	if config.Next == "" {
		*config, err = internal.MakeInitialCall()
		config.Pokedex = savedPokedex
	} else {
		*config, err = internal.GetNextLocation(*config)
	}
	config.SavePath = savedPath
	config.CurrentAreaPokemon = savedArea
	if err != nil {
		return err
	}
	return selectAndExploreLocation(config)
}

func commandMapBack(config *internal.Config, explore string, second string) error {
	var err error
	savedPath := config.SavePath
	savedArea := config.CurrentAreaPokemon
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		*config, err = internal.GetLastLocation(*config)
	}
	config.SavePath = savedPath
	config.CurrentAreaPokemon = savedArea
	if err != nil {
		return err
	}
	return selectAndExploreLocation(config)
}

func selectAndExploreLocation(config *internal.Config) error {
	names := make([]string, 0, len(config.Results))
	for _, loc := range config.Results {
		names = append(names, loc.Name)
	}
	selected, err := pterm.DefaultInteractiveSelect.WithOptions(names).Show()
	if err != nil {
		return err
	}
	return commandExplore(config, selected, "")
}

func commandExplore(config *internal.Config, location string, second string) error {
	pokemon, err := internal.GetPokemonAtLocation(location)
	if err != nil {
		return err
	}

	config.CurrentAreaPokemon = make(map[string]bool, len(pokemon))
	options := make([]string, 0, len(pokemon)+1)
	for _, p := range pokemon {
		config.CurrentAreaPokemon[p.Name] = true
		options = append(options, p.Name)
	}
	options = append(options, "Leave")

	locationHeader := strings.ToUpper(strings.ReplaceAll(location, "-", " "))

	for {
		fmt.Print("\033[H\033[2J")
		pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint(locationHeader))
		fmt.Println()

		selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Catch a Pokemon?")
		if err != nil {
			return err
		}
		if selected == "Leave" {
			return nil
		}

		fmt.Print("\033[H\033[2J")
		pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint(locationHeader))
		fmt.Println()
		if err := commandCatch(config, selected, ""); err != nil {
			return err
		}

		stay, err := pterm.DefaultInteractiveConfirm.WithDefaultText("Stay in area?").Show()
		if err != nil {
			return err
		}
		if !stay {
			return nil
		}
	}
}

func fetchSpriteString(url string, width, height int) string {
	if url == "" {
		return ""
	}
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return ""
	}

	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(img, &convert.Options{
		Colored:     true,
		FixedWidth:  width,
		FixedHeight: height,
	})
}

func renderSprite(url string) {
	art := fetchSpriteString(url, 40, 20)
	if art == "" {
		return
	}
	ansiStrip := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	termWidth := pterm.GetTerminalWidth()
	for _, line := range strings.Split(art, "\n") {
		visible := len(ansiStrip.ReplaceAllString(line, ""))
		pad := (termWidth - visible) / 2
		if pad < 0 {
			pad = 0
		}
		fmt.Println(strings.Repeat(" ", pad) + line)
	}
}

func commandCatch(config *internal.Config, selectedPokemon string, second string) error {
	if len(config.CurrentAreaPokemon) == 0 {
		fmt.Println("You need to explore an area first before you can catch pokemon!")
		return nil
	}
	if !config.CurrentAreaPokemon[selectedPokemon] {
		fmt.Printf("%s is not in this area! Use explore to find pokemon you can catch.\n", selectedPokemon)
		return nil
	}
	userPokemon, err := internal.FindPokemon(selectedPokemon)
	if err != nil {
		return err
	}
	_, exists := config.Pokedex[selectedPokemon]
	if exists {
		fmt.Printf("You already caught %v!\n", selectedPokemon)
		return nil
	}
	
	spinner, _ := pterm.DefaultSpinner.Start("Throwing a Pokeball at " + selectedPokemon + "...")
	time.Sleep(time.Second * 3)
	spinner.Stop()

	userChance := rand.Intn(640)
	if userChance >= userPokemon.BaseExperience {
		config.Pokedex[selectedPokemon] = userPokemon
		spinner.Success(selectedPokemon + " was caught!\n")
		renderSprite(userPokemon.Sprites.FrontDefault)
	} else {
		spinner.Fail(selectedPokemon + " escaped!\n")
	}
	return nil
}

func commandPokedex(config *internal.Config, selectedPokemon string, second string) error {

	if len(config.Pokedex) < 1 {
		fmt.Println("You dont have any pokemon! Get Catching!")
		return nil
	}

	pterm.DefaultCenter.Println(
		pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("YOUR POKEDEX"),
	)

	names := make([]string, 0, len(config.Pokedex))
	for name := range config.Pokedex {
		names = append(names, name)
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(names).Show()
	if err != nil {
		return err
	}

	fmt.Print("\033[H\033[2J")
	return commandInspect(config, selected, "")
}

func commandInspect(config *internal.Config, selectedPokemon string, second string) error {
	if len(config.Pokedex) < 1 {
		fmt.Println("You must catch a pokemon first to inspect it!")
		return nil
	}
	userPokemon, exists := config.Pokedex[selectedPokemon]
	if !exists {
		fmt.Printf("%s doesnt exist in your Pokedex yet\n", selectedPokemon)
		return nil
	}

	renderSprite(userPokemon.Sprites.FrontDefault)

	pterm.DefaultCenter.Println(
		pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint(strings.ToUpper(userPokemon.Name)),
	)

	statsData := pterm.TableData{{"STAT", "VALUE"}}
	for _, stat := range userPokemon.Stats {
		statsData = append(statsData, []string{strings.ToUpper(stat.Stat.Name), fmt.Sprintf("%d", stat.BaseStat)})
	}
	statsTable, _ := pterm.DefaultTable.WithHasHeader().WithData(statsData).Srender()
	pterm.DefaultCenter.Println(statsTable)

	abilitiesData := pterm.TableData{{"ABILITIES"}}
	for _, a := range userPokemon.Abilities {
		abilitiesData = append(abilitiesData, []string{strings.ToUpper(a.Ability.Name)})
	}
	abilitiesTable, _ := pterm.DefaultTable.WithHasHeader().WithData(abilitiesData).Srender()
	pterm.DefaultCenter.Println(abilitiesTable)

	return nil
}

func getDamagingMove(pokemon internal.Pokemon, fallbackAttack int) (string, int) {
	moves := pokemon.Moves
	if len(moves) == 0 {
		fresh, err := internal.FindPokemon(pokemon.Name)
		if err == nil {
			moves = fresh.Moves
		}
	}
	if len(moves) == 0 {
		return "attack", fallbackAttack
	}
	indices := rand.Perm(len(moves))
	if len(indices) > 10 {
		indices = indices[:10]
	}
	for _, i := range indices {
		move, err := internal.GetMove(moves[i].Move.Name)
		if err != nil || move.Power <= 0 {
			continue
		}
		return move.Name, move.Power
	}
	return "attack", fallbackAttack
}

func printHPBar(name string, current, max int) {
	pct := float64(current) / float64(max)
	filled := int(pct * 20)
	empty := 20 - filled
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"

	color := pterm.FgGreen
	if pct < 0.5 {
		color = pterm.FgYellow
	}
	if pct < 0.25 {
		color = pterm.FgRed
	}
	pterm.DefaultCenter.Println(color.Sprintf("%s HP: %s %d/%d", strings.ToUpper(name), bar, current, max))
}

func commandBattle(config *internal.Config, firstPokemon string, secondPokemon string) error {
	pokemonStats, err := getPokemonStats(config, firstPokemon, secondPokemon)
	if err != nil {
		fmt.Printf("Error fetching pokemon stats: %v\n", err)
	}
	pokemonOne := pokemonStats[0]
	pokemonTwo := pokemonStats[1]

	maxHpOne := pokemonOne.hp
	maxHpTwo := pokemonTwo.hp

	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("⚔  BATTLE START  ⚔"))
	fmt.Println()

	spriteOne := fetchSpriteString(config.Pokedex[firstPokemon].Sprites.FrontDefault, 28, 14)
	spriteTwo := fetchSpriteString(config.Pokedex[secondPokemon].Sprites.FrontDefault, 28, 14)
	if spriteOne != "" && spriteTwo != "" {
		panelOut, _ := pterm.DefaultPanel.WithPadding(6).WithPanels(pterm.Panels{
			{{Data: spriteOne}, {Data: spriteTwo}},
		}).Srender()
		ansiStrip := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
		lines := strings.Split(panelOut, "\n")
		maxVisible := 0
		for _, l := range lines {
			if w := len(ansiStrip.ReplaceAllString(l, "")); w > maxVisible {
				maxVisible = w
			}
		}
		pad := strings.Repeat(" ", max((pterm.GetTerminalWidth()-maxVisible)/2, 0))
		for _, l := range lines {
			fmt.Println(pad + l)
		}
		fmt.Println()
	}

	round := 1
	for !pokemonOne.dead || !pokemonTwo.dead {
		pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgYellow, pterm.Bold).Sprintf("— ROUND %d —", round))
		round++

		moveOneName, moveOnePower := getDamagingMove(config.Pokedex[pokemonOne.name], pokemonOne.attack)
		pokemonTwo.hp -= moveOnePower
		if pokemonTwo.hp < 1 {
			pokemonTwo.hp = 0
		}
		pterm.FgYellow.Printf("  %s uses %s on %s for %d damage!\n", strings.ToUpper(pokemonOne.name), strings.ToUpper(moveOneName), strings.ToUpper(pokemonTwo.name), moveOnePower)
		time.Sleep(500 * time.Millisecond)
		printHPBar(pokemonTwo.name, pokemonTwo.hp, maxHpTwo)
		time.Sleep(500 * time.Millisecond)

		if pokemonTwo.hp < 1 {
			pokemonTwo.dead = true
			fmt.Println()
			pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle(strings.ToUpper(pokemonOne.name)+" WINS", pterm.NewStyle(pterm.FgGreen))).Render()
			return nil
		}

		moveTwoName, moveTwoPower := getDamagingMove(config.Pokedex[pokemonTwo.name], pokemonTwo.attack)
		pokemonOne.hp -= moveTwoPower
		if pokemonOne.hp < 1 {
			pokemonOne.hp = 0
		}
		pterm.FgYellow.Printf("  %s uses %s on %s for %d damage!\n", strings.ToUpper(pokemonTwo.name), strings.ToUpper(moveTwoName), strings.ToUpper(pokemonOne.name), moveTwoPower)
		time.Sleep(500 * time.Millisecond)
		printHPBar(pokemonOne.name, pokemonOne.hp, maxHpOne)
		time.Sleep(500 * time.Millisecond)

		if pokemonOne.hp < 1 {
			pokemonOne.dead = true
			fmt.Println()
			pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle(strings.ToUpper(pokemonTwo.name)+" WINS", pterm.NewStyle(pterm.FgGreen))).Render()
			return nil
		}

		fmt.Println()
		time.Sleep(800 * time.Millisecond)
	}
	return nil
}

func commandRelease (config *internal.Config, pokemonToRelease string, second string) error {
	_, exists := config.Pokedex[pokemonToRelease]
	if !exists {
		fmt.Printf("You dont have %v in your pokedex.\n", pokemonToRelease)
		return nil
	}
	delete(config.Pokedex, pokemonToRelease)
	fmt.Printf("%v was released!\n", pokemonToRelease)

	return nil
}

func commandCompare(config *internal.Config, first string, second string) error {
	if len(config.Pokedex) < 2 {
		fmt.Println("You need at least 2 pokemon in your pokedex to compare!")
		return nil
	}

	names := make([]string, 0, len(config.Pokedex))
	for name := range config.Pokedex {
		names = append(names, name)
	}

	firstName := first
	if firstName == "" {
		var err error
		firstName, err = pterm.DefaultInteractiveSelect.WithOptions(names).Show("Select first Pokemon")
		if err != nil {
			return err
		}
	}

	remaining := make([]string, 0, len(names)-1)
	for _, n := range names {
		if n != firstName {
			remaining = append(remaining, n)
		}
	}

	secondName := second
	if secondName == "" {
		var err error
		secondName, err = pterm.DefaultInteractiveSelect.WithOptions(remaining).Show("Select second Pokemon")
		if err != nil {
			return err
		}
	}

	pOne := config.Pokedex[firstName]
	pTwo := config.Pokedex[secondName]

	// index stats by name for comparison
	statsOne := make(map[string]int)
	for _, s := range pOne.Stats {
		statsOne[s.Stat.Name] = s.BaseStat
	}
	statsTwo := make(map[string]int)
	for _, s := range pTwo.Stats {
		statsTwo[s.Stat.Name] = s.BaseStat
	}

	// collect all stat names in order from first pokemon
	statOrder := make([]string, 0, len(pOne.Stats))
	for _, s := range pOne.Stats {
		statOrder = append(statOrder, s.Stat.Name)
	}

	tableOne := pterm.TableData{{"STAT", "VALUE"}}
	tableTwo := pterm.TableData{{"STAT", "VALUE"}}
	for _, statName := range statOrder {
		v1 := statsOne[statName]
		v2 := statsTwo[statName]
		label := strings.ToUpper(statName)

		val1Str := fmt.Sprintf("%d", v1)
		val2Str := fmt.Sprintf("%d", v2)
		if v1 > v2 {
			val1Str = pterm.FgGreen.Sprint(val1Str)
		} else if v2 > v1 {
			val2Str = pterm.FgGreen.Sprint(val2Str)
		}

		tableOne = append(tableOne, []string{label, val1Str})
		tableTwo = append(tableTwo, []string{label, val2Str})
	}

	nameOne := pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint(strings.ToUpper(firstName))
	nameTwo := pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint(strings.ToUpper(secondName))

	spriteOne := fetchSpriteString(pOne.Sprites.FrontDefault, 24, 12)
	spriteTwo := fetchSpriteString(pTwo.Sprites.FrontDefault, 24, 12)

	renderedOne, _ := pterm.DefaultTable.WithHasHeader().WithData(tableOne).Srender()
	renderedTwo, _ := pterm.DefaultTable.WithHasHeader().WithData(tableTwo).Srender()

	buildPanel := func(sprite, name, table string) string {
		if sprite != "" {
			return sprite + "\n" + name + "\n\n" + table
		}
		return name + "\n\n" + table
	}

	panelOutput, _ := pterm.DefaultPanel.WithPadding(6).WithPanels(pterm.Panels{
		{{Data: buildPanel(spriteOne, nameOne, renderedOne)}, {Data: buildPanel(spriteTwo, nameTwo, renderedTwo)}},
	}).Srender()

	ansiStrip := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	lines := strings.Split(panelOutput, "\n")
	maxVisible := 0
	for _, line := range lines {
		if w := len(ansiStrip.ReplaceAllString(line, "")); w > maxVisible {
			maxVisible = w
		}
	}
	leftPad := (pterm.GetTerminalWidth() - maxVisible) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	pad := strings.Repeat(" ", leftPad)

	fmt.Print("\033[H\033[2J")
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("P O K E D E X"))
	fmt.Println()
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgYellow, pterm.Bold).Sprint("— COMPARISON —"))
	fmt.Println()
	for _, line := range lines {
		fmt.Println(pad + line)
	}

	return nil
}
