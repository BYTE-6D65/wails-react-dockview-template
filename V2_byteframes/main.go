package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"V2_byteframes/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Get app data directory
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get user config directory:", err)
	}
	appDataDir := filepath.Join(userConfigDir, "V2_byteframes")

	// Initialize database
	database, err := db.New(appDataDir)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	log.Println("Database initialized at:", database.Path())

	// Create an instance of the app structure
	app := NewApp()
	app.db = database

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "V2_byteframes",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
