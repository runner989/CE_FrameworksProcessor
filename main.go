package main

import (
	"cefp/database"
	"cefp/secret"
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	db, err := database.NewDB("cefp.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	apiKey := getAPIKey()
	if apiKey == "" {
		log.Fatal("No API key provided")
	}

	// Create an instance of the app structure with initialized fields
	app := NewApp(apiKey, db)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Compliance Essentials Frameworks Processor",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getAPIKey() string {
	var apiKey string
	// Try to load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using embedded API key.")
		apiKey = secret.APIKey
	}

	// Check if the API key is set in the environment
	if envKey := os.Getenv("API_KEY"); envKey != "" {
		apiKey = envKey
	}
	return apiKey
}
