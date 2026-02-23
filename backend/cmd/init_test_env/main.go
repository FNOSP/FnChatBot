package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	baseURL := "http://localhost:8080/api"

	// Ensure DB is reset or clean (optional, but good for test)
	// We just overwrite the default model config

	config := map[string]interface{}{
		"name":       "GPT-4o Mini CA",
		"provider":   "openai",
		"base_url":   os.Getenv("TEST_BASE_URL"),
		"api_key":    os.Getenv("TEST_API_KEY"),
		"model":      os.Getenv("TEST_MODEL_ID"),
		"is_default": true,
	}

	if config["base_url"] == "" || config["api_key"] == "" {
		fmt.Println("Error: TEST_BASE_URL and TEST_API_KEY env vars required")
		os.Exit(1)
	}

	body, _ := json.Marshal(config)
	resp, err := http.Post(baseURL+"/models", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to configure model: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Printf("Failed to configure model, status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Model configured successfully for Real AI Test")
}
