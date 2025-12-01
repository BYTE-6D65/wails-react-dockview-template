package models

// Setting represents a key-value configuration setting
type Setting struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Layout represents a saved Dockview layout configuration
type Layout struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LayoutJSON string `json:"layout_json"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// WindowState represents the window position and size
type WindowState struct {
	ID        int   `json:"id"`
	X         int   `json:"x"`
	Y         int   `json:"y"`
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	Maximized bool  `json:"maximized"`
	UpdatedAt int64 `json:"updated_at"`
}
