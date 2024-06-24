package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	ErrAreaName         = errors.New("pokeapi: area's cannot be empty")
	ErrGetPokeByAreaApi = errors.New("pokeapi: failed to fetch pokemons")
	ErrAreaNotFound     = errors.New("pokeapi: area not found")
	ErrPokemonNotFound  = errors.New("pokeapi: pokemon not found")
)

var baseUrl = "https://pokeapi.co/api/v2"

type cacheService interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

type PokeApi interface {
	GetLocationAreas(apiUrl string) (GetLocationAreaResp, error)
	GetPokemonsByArea(area string) (GetPokemonAreaByAreaResp, error)
	GetPokemon(pokemon string) (Pokemon, error)
}

type pokeApi struct {
	cache cacheService
}

func NewPokeApi(cs cacheService) *pokeApi {
	return &pokeApi{
		cache: cs,
	}
}

func (p *pokeApi) responseFromCache(key string, v any) (bool, error) {
	item, exist := p.cache.Get(key)

	err := json.Unmarshal(item, v)

	if err != nil {
		return exist, fmt.Errorf("pokeapi: %w", err)
	}

	return exist, nil
}

func (p *pokeApi) GetLocationAreas(apiUrl string) (GetLocationAreaResp, error) {
	var (
		resp GetLocationAreaResp
	)

	if ok, err := p.responseFromCache(apiUrl, &resp); ok && err == nil {
		return resp, nil
	}

	res, err := http.Get(apiUrl)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var v GetLocationAreaResp
	json.Unmarshal(data, &v)

	p.cache.Add(apiUrl, data)
	return v, nil
}

type GetPokemByAreaIn struct {
	AreaName string
}

func (p *pokeApi) GetPokemonsByArea(area string) (GetPokemonAreaByAreaResp, error) {

	var (
		resp GetPokemonAreaByAreaResp
	)

	if ok, err := p.responseFromCache(area, &resp); ok && err == nil {
		return resp, nil
	}

	if area == "" {
		return resp, ErrAreaName
	}

	apiUrl := fmt.Sprintf("%s/location-area/%s/", baseUrl, area)

	r, err := http.Get(apiUrl)

	if err != nil {
		return resp, errors.New("failed to fetch pokemons by area")
	}

	if r.StatusCode == http.StatusNotFound {
		return resp, ErrAreaNotFound
	}

	body, _ := io.ReadAll(r.Body)

	err = json.Unmarshal(body, &resp)

	if err != nil {
		return resp, fmt.Errorf("pokeapi: %w", err)
	}

	p.cache.Add(area, body)

	return resp, nil
}

func (p *pokeApi) GetPokemon(pokemon string) (Pokemon, error) {
	r, err := http.Get(
		fmt.Sprintf("%s/pokemon/%s", baseUrl, pokemon),
	)

	if err != nil {
		return Pokemon{}, fmt.Errorf("pokeapi: %w", err)
	}

	if r.StatusCode == http.StatusNotFound {
		return Pokemon{}, ErrPokemonNotFound
	}

	if r.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("pokeapi: %w", err)
	}

	var poke Pokemon
	err = json.NewDecoder(r.Body).Decode(&poke)

	if err != nil {
		return Pokemon{}, fmt.Errorf("pokeapi: %w", err)
	}

	return poke, nil
}
