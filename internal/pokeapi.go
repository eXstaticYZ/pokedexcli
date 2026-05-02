package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pokedexcli/internal/pokecache"
)

// LocationArea represents a Pokémon location area from the API
type LocationArea struct {
	Name string `json:"name"`
}

// PokeAPIResponse represents the paginated response from the Pokémon API
type PokeAPIResponse struct {
	Results []LocationArea `json:"results"`
	Next    string         `json:"next"`
	Prev    string         `json:"previous"`
}

// Pokemon represents a Pokémon from the Pokémon API
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

// PokeAPIClient provides methods for interacting with the Pokémon API
type PokeAPIClient interface {
	GetLocationAreas(url string) (*PokeAPIResponse, error)
	GetLocationAreaPokemon(url string) ([]Encounter, error)
	GetPokemon(url string) (*Pokemon, error)
}

// RealPokeAPIClient implements the PokeAPIClient interface
type RealPokeAPIClient struct {
	cache *pokecache.Cache
}

// NewRealPokeAPIClient creates a new real Pokémon API client with caching
func NewRealPokeAPIClient(cache *pokecache.Cache) *RealPokeAPIClient {
	return &RealPokeAPIClient{
		cache: cache,
	}
}

// GetLocationAreas fetches location areas from the Pokémon API, using cache when available
func (c *RealPokeAPIClient) GetLocationAreas(url string) (*PokeAPIResponse, error) {
	// First check if we have this URL in cache
	if c.cache != nil {
		cachedData, found := c.cache.Get(url)
		if found {
			var apiResponse PokeAPIResponse
			err := json.Unmarshal(cachedData, &apiResponse)
			if err == nil {
				fmt.Printf("[CACHE HIT] Using cached data for: %s\n", url)
				return &apiResponse, nil
			} else {
				fmt.Printf("[CACHE ERROR] Cache hit but invalid data for: %s\n", url)
			}
		} else {
			fmt.Printf("[CACHE MISS] Making API request for: %s\n", url)
		}
	}

	// If not in cache or cache miss, make the actual API call
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var apiResponse PokeAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Cache the response if we have a cache
	if c.cache != nil {
		cacheData, err := json.Marshal(apiResponse)
		if err == nil {
			c.cache.Add(url, cacheData)
			fmt.Printf("[CACHE STORE] Stored response for: %s\n", url)
		} else {
			fmt.Printf("[CACHE ERROR] Failed to store response for: %s\n", url)
		}
	}

	return &apiResponse, nil
}

// Encounter represents a Pokémon encounter in a location area
type Encounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

// GetLocationAreaPokemon fetches Pokémon encounters for a location area from the Pokémon API
func (c *RealPokeAPIClient) GetLocationAreaPokemon(url string) ([]Encounter, error) {
	// First check if we have this URL in cache
	if c.cache != nil {
		cachedData, found := c.cache.Get(url)
		if found {
			var encounters []Encounter
			err := json.Unmarshal(cachedData, &encounters)
			if err == nil {
				fmt.Printf("[CACHE HIT] Using cached data for: %s\n", url)
				return encounters, nil
			} else {
				fmt.Printf("[CACHE ERROR] Cache hit but invalid data for: %s\n", url)
			}
		} else {
			fmt.Printf("[CACHE MISS] Making API request for: %s\n", url)
		}
	}

	// If not in cache or cache miss, make the actual API call
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch encounters: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var encounters []Encounter
	err = json.Unmarshal(body, &encounters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Cache the response if we have a cache
	if c.cache != nil {
		cacheData, err := json.Marshal(encounters)
		if err == nil {
			c.cache.Add(url, cacheData)
			fmt.Printf("[CACHE STORE] Stored response for: %s\n", url)
		} else {
			fmt.Printf("[CACHE ERROR] Failed to store response for: %s\n", url)
		}
	}

	return encounters, nil
}

// GetPokemon fetches Pokémon information from the Pokémon API
func (c *RealPokeAPIClient) GetPokemon(url string) (*Pokemon, error) {
	// First check if we have this URL in cache
	if c.cache != nil {
		cachedData, found := c.cache.Get(url)
		if found {
			var pokemon Pokemon
			err := json.Unmarshal(cachedData, &pokemon)
			if err == nil {
				fmt.Printf("[CACHE HIT] Using cached data for: %s\n", url)
				return &pokemon, nil
			} else {
				fmt.Printf("[CACHE ERROR] Cache hit but invalid data for: %s\n", url)
			}
		} else {
			fmt.Printf("[CACHE MISS] Making API request for: %s\n", url)
		}
	}

	// If not in cache or cache miss, make the actual API call
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Cache the response if we have a cache
	if c.cache != nil {
		cacheData, err := json.Marshal(pokemon)
		if err == nil {
			c.cache.Add(url, cacheData)
			fmt.Printf("[CACHE STORE] Stored response for: %s\n", url)
		} else {
			fmt.Printf("[CACHE ERROR] Failed to store response for: %s\n", url)
		}
	}

	return &pokemon, nil
}
