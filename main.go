package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"pokedexcli/internal"
	"pokedexcli/internal/pokecache"
)

const cacheInterval = 5 * time.Second

// cliConfig holds the state for pagination
type cliConfig struct {
	nextURL    string
	prevURL    string
	currentURL string
}

// cliCommand represents a command in REPL
type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig, []string) error
}

// global config instance
var globalConfig = &cliConfig{}

// pokeAPIClient is the global API client instance
var pokeAPIClient internal.PokeAPIClient

// cacheInstance is the global cache instance
var cacheInstance *pokecache.Cache

// caughtPokemon is a map to track caught Pokemon
var caughtPokemon = make(map[string]internal.Pokemon)

// commandExit handles the exit command
func commandExit(config *cliConfig, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// commandHelp handles the help command
func commandHelp(config *cliConfig, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Printf("help: %s\n", getCommand("help").description)
	fmt.Printf("exit: %s\n", getCommand("exit").description)
	fmt.Printf("map: %s\n", getCommand("map").description)
	fmt.Printf("mapb: %s\n", getCommand("mapb").description)
	fmt.Printf("catch <pokemon_name>: %s\n", getCommand("catch").description)
	fmt.Printf("caught: %s\n", getCommand("caught").description)
	fmt.Printf("inspect <pokemon_name>: %s\n", getCommand("inspect").description)
	fmt.Printf("pokedex: %s\n", getCommand("pokedex").description)
	fmt.Printf("explore <area_name>: %s\n", getCommand("explore").description)
	return nil
}

// commandMap handles the map command to show location areas
func commandMap(config *cliConfig, params []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.nextURL != "" {
		url = config.nextURL
	}

	apiResponse, err := pokeAPIClient.GetLocationAreas(url)
	if err != nil {
		return err
	}

	fmt.Printf("Displaying locations from: %s\n", url)
	for i, result := range apiResponse.Results {
		cleanedName := cleanInput(result.Name)[0]
		fmt.Printf("%d. %s\n", i+1, cleanedName)
	}
	fmt.Printf("Total: %d locations\n", len(apiResponse.Results))

	config.currentURL = url
	config.nextURL = apiResponse.Next
	config.prevURL = apiResponse.Prev

	return nil
}

// commandMapBack handles going back in map pagination
func commandMapBack(config *cliConfig, params []string) error {
	if config.prevURL == "" {
		return fmt.Errorf("no previous page available")
	}

	url := config.prevURL

	apiResponse, err := pokeAPIClient.GetLocationAreas(url)
	if err != nil {
		return err
	}

	// Display location names
	fmt.Printf("Displaying previous locations from: %s\n", url)
	for i, result := range apiResponse.Results {
		cleanedName := cleanInput(result.Name)[0]
		fmt.Printf("%d. %s\n", i+1, cleanedName)
	}
	fmt.Printf("Total: %d locations\n", len(apiResponse.Results))

	// Update config for next/previous navigation
	config.currentURL = url
	config.nextURL = apiResponse.Next
	config.prevURL = apiResponse.Prev

	return nil
}

// commandExplore handles the explore command to show Pokemon in a location area
func commandExplore(config *cliConfig, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("please specify a location area name")
	}

	areaName := params[0]
	fmt.Printf("Exploring %s...\n", areaName)

	// Construct the URL for Pokemon encounters
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)

	encounters, err := pokeAPIClient.GetLocationAreaPokemon(url)
	if err != nil {
		return err
	}

	// Display Pokemon names
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounters {
		cleanedName := cleanInput(encounter.Pokemon.Name)[0]
		fmt.Printf(" - %s\n", cleanedName)
	}
	fmt.Printf("Total: %d Pokemon found\n", len(encounters))

	return nil
}

// getCommand looks up a command in the
// commandCatch handles the catch command
func commandCatch(config *cliConfig, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("please specify a Pokemon name")
	}

	pokemonName := params[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Construct the URL for the Pokemon
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	pokemon, err := pokeAPIClient.GetPokemon(url)
	if err != nil {
		return err
	}

	// Calculate catch probability based on base_experience
	// Higher base_experience means harder to catch (inverse relationship)
	maxBaseExp := 300.0 // Maximum reasonable base experience value
	catchProbability := maxBaseExp / (float64(pokemon.BaseExperience) + maxBaseExp)
	catchChance := rand.Float64()

	if catchChance < catchProbability {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		caughtPokemon[pokemon.Name] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

// commandCaught handles the caught command to show all caught Pokemon
func commandCaught(config *cliConfig, params []string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name, pokemon := range caughtPokemon {
		fmt.Printf("- %s (Base Exp: %d)\n", name, pokemon.BaseExperience)
	}
	fmt.Printf("Total caught: %d\n", len(caughtPokemon))
	return nil
}

// commandInspect handles the inspect command to show details of a caught Pokemon
func commandInspect(config *cliConfig, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("please specify a Pokemon name")
	}

	pokemonName := params[0]
	pokemon, exists := caughtPokemon[pokemonName]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	// Display stats from the caught Pokemon
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	// Display types
	fmt.Println("Types:")
	for i, typeInfo := range pokemon.Types {
		if i > 0 {
			fmt.Printf("  - %s\n", typeInfo.Type.Name)
		} else {
			fmt.Printf("  - %s\n", typeInfo.Type.Name)
		}
	}

	return nil
}

// commandPokedex handles the pokedex command to show names of all caught Pokemon
func commandPokedex(config *cliConfig, params []string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for i, name := range getSortedPokemonNames() {
		fmt.Printf("%d. %s\n", i+1, name)
	}
	fmt.Printf("Total: %d Pokemon\n", len(caughtPokemon))
	return nil
}

// getSortedPokemonNames returns a sorted list of caught Pokemon names
func getSortedPokemonNames() []string {
	names := make([]string, 0, len(caughtPokemon))
	for name := range caughtPokemon {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// getCommand looks up a command in the registry by name
func getCommand(name string) cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows location areas (20 at a time)",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Goes back to previous set of locations",
			callback:    commandMapBack,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback:    commandCatch,
		},
		"caught": {
			name:        "caught",
			description: "Shows all caught Pokemon",
			callback:    commandCaught,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon's details",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows names of all caught Pokemon",
			callback:    commandPokedex,
		},
		"explore": {
			name:        "explore",
			description: "Shows Pokemon found in a specific location area",
			callback:    commandExplore,
		},
	}
	return commands[name]
}

func main() {
	// Initialize the cache
	cacheInstance = pokecache.NewCache(cacheInterval)
	defer cacheInstance.Close()

	// Initialize the PokeAPI client with caching
	pokeAPIClient = internal.NewRealPokeAPIClient(cacheInstance)

	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Start an infinite loop for the REPL
	for {
		fmt.Print("Pokedex > ")

		// Check if we can scan input
		if !scanner.Scan() {
			break
		}

		// Get the user's input
		input := scanner.Text()

		// Clean the input
		cleanedInput := cleanInput(input)

		// Check if there's any input
		if len(cleanedInput) > 0 {
			commandName := cleanedInput[0]
			parameters := []string{}

			// If there are additional words, treat them as parameters
			if len(cleanedInput) > 1 {
				parameters = cleanedInput[1:]
			}

			// Look up the command in registry
			command := getCommand(commandName)
			if command.callback != nil {
				if err := command.callback(globalConfig, parameters); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
