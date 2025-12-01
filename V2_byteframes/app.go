package main

import (
	"context"
	"fmt"

	"V2_byteframes/db"
	"V2_byteframes/models"
)

// App struct
type App struct {
	ctx context.Context
	db  *db.Database
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is shutting down
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// --- Settings Methods ---

// GetSetting retrieves a setting by key
func (a *App) GetSetting(key string) (string, error) {
	return a.db.GetSetting(key)
}

// SetSetting creates or updates a setting
func (a *App) SetSetting(key, value string) error {
	return a.db.SetSetting(key, value)
}

// GetAllSettings retrieves all settings
func (a *App) GetAllSettings() ([]models.Setting, error) {
	return a.db.GetAllSettings()
}

// --- Layout Methods ---

// SaveLayout creates or updates a layout
func (a *App) SaveLayout(name, layoutJSON string) (*models.Layout, error) {
	return a.db.SaveLayout(name, layoutJSON)
}

// GetAllLayouts retrieves all saved layouts
func (a *App) GetAllLayouts() ([]models.Layout, error) {
	return a.db.GetAllLayouts()
}

// GetLayout retrieves a layout by ID
func (a *App) GetLayout(id int) (*models.Layout, error) {
	return a.db.GetLayout(id)
}

// SetActiveLayout sets a layout as active
func (a *App) SetActiveLayout(id int) error {
	return a.db.SetActiveLayout(id)
}

// GetActiveLayout retrieves the currently active layout
func (a *App) GetActiveLayout() (*models.Layout, error) {
	return a.db.GetActiveLayout()
}

// DeleteLayout deletes a layout by ID
func (a *App) DeleteLayout(id int) error {
	return a.db.DeleteLayout(id)
}

// --- Window State Methods ---

// SaveWindowState saves the window position and size
func (a *App) SaveWindowState(x, y, width, height int, maximized bool) error {
	return a.db.SaveWindowState(x, y, width, height, maximized)
}

// GetWindowState retrieves the saved window state
func (a *App) GetWindowState() (*models.WindowState, error) {
	return a.db.GetWindowState()
}
