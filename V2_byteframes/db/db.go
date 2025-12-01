package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Database holds the SQLite connection and metadata
type Database struct {
	DB      *sql.DB
	path    string
	version int
}

// New creates a new database connection
func New(appDataDir string) (*Database, error) {
	// Ensure app data directory exists
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app data directory: %w", err)
	}

	dbPath := filepath.Join(appDataDir, "app.db")

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)

	database := &Database{
		DB:      db,
		path:    dbPath,
		version: 1,
	}

	// Run migrations
	if err := database.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Path returns the database file path
func (d *Database) Path() string {
	return d.path
}
