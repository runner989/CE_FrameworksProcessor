package main

import (
	"cefp/database"
	"cefp/secret"
	"database/sql"
	"embed"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	memDB, err := database.CreateInMemoryDB()
	if err != nil {
		log.Fatal(err)
	}
	err = database.InitializeMemoryDB(memDB)
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDB("cefp.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Failed to close database connection: %v", err)
		}
	}(db)

	apiKey := getAPIKey()
	if apiKey == "" {
		log.Fatal("No API key provided")
	}

	// Create an instance of the app structure with initialized fields
	app := NewApp(apiKey, db, memDB)

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "Compliance Essentials Frameworks Processor",
		Width:            1020,
		Height:           768,
		Fullscreen:       false,
		DisableResize:    false,
		WindowStartState: options.Normal,
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

// getAPIKey check .env to see if new API Key is there
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
	} else {
		apiKey = secret.APIKey
	}
	return apiKey
}
