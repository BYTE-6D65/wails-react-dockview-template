package db

import (
	"database/sql"
	"time"

	"V2_byteframes/models"
)

// --- Settings Queries ---

// GetSetting retrieves a setting by key
func (d *Database) GetSetting(key string) (string, error) {
	var value string
	err := d.DB.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Return empty string if not found
		}
		return "", err
	}
	return value, nil
}

// SetSetting creates or updates a setting
func (d *Database) SetSetting(key, value string) error {
	_, err := d.DB.Exec(`
		INSERT INTO settings (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`, key, value, value)
	return err
}

// GetAllSettings retrieves all settings
func (d *Database) GetAllSettings() ([]models.Setting, error) {
	rows, err := d.DB.Query(`
		SELECT key, value,
		       strftime('%s', created_at) as created_at,
		       strftime('%s', updated_at) as updated_at
		FROM settings
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var s models.Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, rows.Err()
}

// --- Layout Queries ---

// SaveLayout creates or updates a layout
func (d *Database) SaveLayout(name, layoutJSON string) (*models.Layout, error) {
	// Check if layout with this name exists
	var existingID int
	err := d.DB.QueryRow("SELECT id FROM layouts WHERE name = ?", name).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new layout
		result, err := d.DB.Exec(`
			INSERT INTO layouts (name, layout_json)
			VALUES (?, ?)
		`, name, layoutJSON)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		return d.GetLayout(int(id))
	} else if err != nil {
		return nil, err
	}

	// Update existing layout
	_, err = d.DB.Exec(`
		UPDATE layouts
		SET layout_json = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, layoutJSON, existingID)
	if err != nil {
		return nil, err
	}

	return d.GetLayout(existingID)
}

// GetLayout retrieves a layout by ID
func (d *Database) GetLayout(id int) (*models.Layout, error) {
	var l models.Layout
	err := d.DB.QueryRow(`
		SELECT id, name, layout_json, is_active,
		       strftime('%s', created_at) as created_at,
		       strftime('%s', updated_at) as updated_at
		FROM layouts
		WHERE id = ?
	`, id).Scan(&l.ID, &l.Name, &l.LayoutJSON, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// GetAllLayouts retrieves all saved layouts
func (d *Database) GetAllLayouts() ([]models.Layout, error) {
	rows, err := d.DB.Query(`
		SELECT id, name, layout_json, is_active,
		       strftime('%s', created_at) as created_at,
		       strftime('%s', updated_at) as updated_at
		FROM layouts
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layouts []models.Layout
	for rows.Next() {
		var l models.Layout
		if err := rows.Scan(&l.ID, &l.Name, &l.LayoutJSON, &l.IsActive, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		layouts = append(layouts, l)
	}
	return layouts, rows.Err()
}

// SetActiveLayout sets a layout as the active one
func (d *Database) SetActiveLayout(id int) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deactivate all layouts
	if _, err := tx.Exec("UPDATE layouts SET is_active = 0"); err != nil {
		return err
	}

	// Activate the specified layout
	if _, err := tx.Exec("UPDATE layouts SET is_active = 1 WHERE id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}

// GetActiveLayout retrieves the currently active layout
func (d *Database) GetActiveLayout() (*models.Layout, error) {
	var l models.Layout
	err := d.DB.QueryRow(`
		SELECT id, name, layout_json, is_active,
		       strftime('%s', created_at) as created_at,
		       strftime('%s', updated_at) as updated_at
		FROM layouts
		WHERE is_active = 1
		LIMIT 1
	`).Scan(&l.ID, &l.Name, &l.LayoutJSON, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active layout
		}
		return nil, err
	}
	return &l, nil
}

// DeleteLayout deletes a layout by ID
func (d *Database) DeleteLayout(id int) error {
	_, err := d.DB.Exec("DELETE FROM layouts WHERE id = ?", id)
	return err
}

// --- Window State Queries ---

// SaveWindowState saves the window position and size
func (d *Database) SaveWindowState(x, y, width, height int, maximized bool) error {
	_, err := d.DB.Exec(`
		INSERT INTO window_state (id, x, y, width, height, maximized, updated_at)
		VALUES (1, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			x = ?,
			y = ?,
			width = ?,
			height = ?,
			maximized = ?,
			updated_at = ?
	`, x, y, width, height, maximized, time.Now().Unix(),
		x, y, width, height, maximized, time.Now().Unix())
	return err
}

// GetWindowState retrieves the saved window state
func (d *Database) GetWindowState() (*models.WindowState, error) {
	var ws models.WindowState
	err := d.DB.QueryRow(`
		SELECT id, x, y, width, height, maximized,
		       strftime('%s', updated_at) as updated_at
		FROM window_state
		WHERE id = 1
	`).Scan(&ws.ID, &ws.X, &ws.Y, &ws.Width, &ws.Height, &ws.Maximized, &ws.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No saved state
		}
		return nil, err
	}
	return &ws, nil
}
