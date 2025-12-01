package db

import (
	"fmt"
)

// migrate runs all database migrations
func (d *Database) migrate() error {
	// Create schema_version table if it doesn't exist
	if err := d.createVersionTable(); err != nil {
		return err
	}

	// Get current version
	currentVersion, err := d.getCurrentVersion()
	if err != nil {
		return err
	}

	// Run migrations in order
	migrations := []migration{
		{version: 1, up: d.migration001_initial},
	}

	for _, m := range migrations {
		if m.version > currentVersion {
			if err := m.up(); err != nil {
				return fmt.Errorf("migration %d failed: %w", m.version, err)
			}
			if err := d.setVersion(m.version); err != nil {
				return fmt.Errorf("failed to update version to %d: %w", m.version, err)
			}
			fmt.Printf("Applied migration %d\n", m.version)
		}
	}

	return nil
}

type migration struct {
	version int
	up      func() error
}

// createVersionTable creates the schema version tracking table
func (d *Database) createVersionTable() error {
	_, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			version INTEGER NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// getCurrentVersion returns the current schema version
func (d *Database) getCurrentVersion() (int, error) {
	var version int
	err := d.DB.QueryRow("SELECT version FROM schema_version WHERE id = 1").Scan(&version)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

// setVersion updates the schema version
func (d *Database) setVersion(version int) error {
	_, err := d.DB.Exec(`
		INSERT INTO schema_version (id, version)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET version = ?, updated_at = CURRENT_TIMESTAMP
	`, version, version)
	return err
}

// migration001_initial creates the initial schema
func (d *Database) migration001_initial() error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Settings table - for app preferences
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Layouts table - for saved Dockview layouts
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS layouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			layout_json TEXT NOT NULL,
			is_active INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Window state table - for window position and size
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS window_state (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			x INTEGER,
			y INTEGER,
			width INTEGER,
			height INTEGER,
			maximized INTEGER DEFAULT 0,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create index on layouts.is_active for faster lookup
	_, err = tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_layouts_active ON layouts(is_active)
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
