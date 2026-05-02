package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pokedexcli/internal/pokecache"
)

func TestRealPokeAPIClient_GetLocationAreas(t *testing.T) {
	// Create a test server with mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"results": [
				{"name": "test-location-1"},
				{"name": "test-location-2"}
			],
			"next": "https://api.example.com/next",
			"previous": null
		}`))
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Second)
	defer cache.Close()
	client := NewRealPokeAPIClient(cache)

	response, err := client.GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(response.Results))
	}

	if response.Results[0].Name != "test-location-1" {
		t.Errorf("Expected first location to be 'test-location-1', got '%s'", response.Results[0].Name)
	}

	if response.Next != "https://api.example.com/next" {
		t.Errorf("Expected next URL, got '%s'", response.Next)
	}
}

func TestRealPokeAPIClient_GetLocationAreas_EmptyResponse(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	defer cache.Close()
	client := NewRealPokeAPIClient(cache)

	// Create a test server with empty response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"results": [],
			"next": null,
			"previous": null
		}`))
	}))
	defer server.Close()

	response, err := client.GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(response.Results))
	}
}

func TestRealPokeAPIClient_GetLocationAreas_InvalidJSON(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	defer cache.Close()
	client := NewRealPokeAPIClient(cache)

	// Create a test server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	response, err := client.GetLocationAreas(server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON, got none")
	}

	if response != nil {
		t.Errorf("Expected nil response for error case, got %v", response)
	}
}

func TestRealPokeAPIClient_GetLocationAreas_Caching(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	defer cache.Close()
	client := NewRealPokeAPIClient(cache)

	// Create a test server with mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"results": [
				{"name": "cached-location"}
			],
			"next": null,
			"previous": null
		}`))
	}))
	defer server.Close()

	// First call - should hit the API
	response1, err := client.GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	if len(response1.Results) != 1 {
		t.Errorf("Expected 1 result on first call, got %d", len(response1.Results))
	}

	// Second call - should use cache
	response2, err := client.GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if len(response2.Results) != 1 {
		t.Errorf("Expected 1 result on second call, got %d", len(response2.Results))
	}
}
