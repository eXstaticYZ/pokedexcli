package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Start an infinite loop for the REPL
	for {
		// Print prompt without newline
		fmt.Print("Pokedex > ")

		// Check if we can scan input
		if !scanner.Scan() {
			// If scanning fails (e.g., EOF), exit the program
			break
		}

		// Get the user's input
		input := scanner.Text()

		// Clean the input using our function
		cleanedInput := cleanInput(input)

		// Capture the first word if there is any input
		if len(cleanedInput) > 0 {
			firstWord := cleanedInput[0]
			fmt.Printf("Your command was: %s\n", firstWord)
		}
	}

	// Check for scanner errors (e.g., read errors)
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
