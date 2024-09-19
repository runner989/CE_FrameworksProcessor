package main

import (
	"cefp/database"
	"cefp/secret"
	"database/sql"
	"embed"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type apiConfig struct {
	mu     sync.Mutex
	db     *sql.DB
	apiKey string
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	db, err := database.NewDB("cefp.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	cfg := &apiConfig{
		apiKey: getAPIKey(),
		db:     db,
	}
	if cfg.apiKey == "" {
		log.Fatal("no API key provided")
	}

	// _ = &apiConfig{db: db, apiKey: apiKey}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Compliance Essentials Frameworks Processor",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
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
