package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/enjaytarigan/pokedexcli/internal/pokeapi"
)

type cliConfig struct {
	name              string
	nextLocAreaApiUrl string
	prevLocAreaApiUrl string
	pokeApi           pokeapi.PokeApi
	caughPokemon      map[string]pokeapi.Pokemon
}

func (c *cliConfig) AddCaughtPokemon(pokemon pokeapi.Pokemon) {
	c.caughPokemon[pokemon.Name] = pokemon
}

func (c cliConfig) GetCaughtPokemon(pokemonName string) (pokemon pokeapi.Pokemon, exist bool) {
	poke, ok := c.caughPokemon[pokemonName]

	return poke, ok
}

func (c cliConfig) GetAllCaughtPokemons() map[string]pokeapi.Pokemon {
	return c.caughPokemon
}

func (c cliConfig) NextLocAreaApiUrl() string {
	if c.nextLocAreaApiUrl == "" {
		return "https://pokeapi.co/api/v2/location-area/"
	}

	return c.nextLocAreaApiUrl
}

func (c *cliConfig) SetLocAreaApiUrl(nextUrl string, prevUrl string) {
	c.nextLocAreaApiUrl = nextUrl
	c.prevLocAreaApiUrl = prevUrl
}

type cliCommand struct {
	name        string
	description string
	handler     func(cfg *cliConfig, args ...string) error
}

func startRepl(cli *cliConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s > ", cli.name)
		scanner.Scan()

		words := clearInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		var args []string
		if len(words) > 1 {
			args = words[1:]
		}

		commandName := words[0]
		if cmd, ok := getCommands()[commandName]; ok {
			cmd.handler(cli, args...)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func clearInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			handler:     handleCommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			handler:     handleCommandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 location areas",
			handler:     handleCommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			handler:     handleCommandMapb,
		},
		"explore": {
			name:        "explore <area's name>",
			description: "Display list of Pokemons in a given area",
			handler:     handleCommandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Catch a pokemon",
			handler:     handleCatchPokemon,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Inpsect a caught pokemon",
			handler:     handleInspectPokemon,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display pokedex",
			handler:     handleCommandPokedex,
		},
	}
}
