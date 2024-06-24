package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/enjaytarigan/pokedexcli/internal/pokeapi"
)

func handleCommandHelp(cfg *cliConfig, args ...string) error {
	fmt.Printf("Welcome to the %s!\n", cfg.name)
	fmt.Println("Usage:")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func handleCommandExit(cfg *cliConfig, args ...string) error {
	os.Exit(0)
	return nil
}

func handleCommandMap(cfg *cliConfig, args ...string) error {
	locArea, err := cfg.pokeApi.GetLocationAreas(cfg.NextLocAreaApiUrl())

	if err != nil {
		log.Fatal(err)
	}

	cfg.SetLocAreaApiUrl(locArea.Next, locArea.Prev)
	for i, area := range locArea.Results {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}

	return nil
}

func handleCommandMapb(cfg *cliConfig, args ...string) error {
	if cfg.prevLocAreaApiUrl == "" {
		fmt.Println("there are no previous locations")
		return nil
	}

	locArea, err := cfg.pokeApi.GetLocationAreas(cfg.prevLocAreaApiUrl)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(locArea.Next, locArea.Prev)
	cfg.SetLocAreaApiUrl(locArea.Next, locArea.Prev)
	for i, area := range locArea.Results {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}

	return nil
}

func handleCommandExplore(cfg *cliConfig, args ...string) error {
	if len(args) == 0 {
		fmt.Println("you must provide location area")
		return nil
	}

	area := args[0]
	fmt.Printf("Exploring %s...\n", area)
	resp, err := cfg.pokeApi.GetPokemonsByArea(area)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Found pokemon: ")
	for _, p := range resp.PokemonEncounters {
		fmt.Printf("- %s\n", p.Pokemon.Name)
	}

	return nil
}

func handleCatchPokemon(cfg *cliConfig, args ...string) error {
	if len(args) == 0 {
		fmt.Println("you must provide a pokemon name")
		return nil
	}

	pokemonName := args[0]
	pokemon, err := cfg.pokeApi.GetPokemon(pokemonName)

	if errors.Is(err, pokeapi.ErrPokemonNotFound) {
		fmt.Printf("%s pokemon not found\n", pokemonName)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	res := rand.Intn(pokemon.BaseExperience)

	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	cfg.AddCaughtPokemon(pokemon)
	fmt.Printf("%s was caught!\n", pokemon.Name)
	return nil
}

func handleInspectPokemon(cfg *cliConfig, args ...string) error {
	if len(args) == 0 {
		fmt.Println("you must provide a pokemon name")
		return nil
	}

	pokemon, exist := cfg.GetCaughtPokemon(args[0])

	if !exist {
		fmt.Println("Pokemon is not caught yet")
		return nil
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)

	fmt.Println("Stats: ")
	for _, s := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types: ")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func handleCommandPokedex(cfg *cliConfig, args ...string) error {
	fmt.Println("Your Pokedex: ")
	for _, pokemon := range cfg.GetAllCaughtPokemons() {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}
