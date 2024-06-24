package main

import (
	"time"

	"github.com/enjaytarigan/pokedexcli/internal/pokeapi"
	"github.com/enjaytarigan/pokedexcli/internal/pokecache"
)

func main() {
	cacheService := pokecache.NewCache(5 * time.Second)

	cfg := &cliConfig{
		pokeApi:      pokeapi.NewPokeApi(cacheService),
		name:         "Pokedex",
		caughPokemon: make(map[string]pokeapi.Pokemon),
	}
	startRepl(cfg)
}
