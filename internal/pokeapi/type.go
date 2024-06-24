package pokeapi

type LocationArea struct {
	Name string `json:"name"`
}

type PokemoType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Pokemon struct {
	Name           string       `json:"name"`
	BaseExperience int          `json:"base_experience,omitempty"`
	Height         int          `json:"height"`
	Weight         int          `json:"weight"`
	Types          []PokemoType `json:"types"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type GetPokemonAreaByAreaResp struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type GetLocationAreaResp struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Prev    string         `json:"previous"`
	Results []LocationArea `json:"results"`
}
