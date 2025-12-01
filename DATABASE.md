# SQLite Database Integration

This template includes a complete SQLite database integration using **modernc.org/sqlite** (pure Go, no CGO required).

## Features

- ✅ **Pure Go SQLite** - No CGO dependencies
- ✅ **Migration System** - Version-based schema management
- ✅ **Layout Persistence** - Save and load Dockview layouts
- ✅ **Theme Persistence** - Remember user's theme choice
- ✅ **Settings Storage** - Key-value configuration store
- ✅ **Window State** - Save/restore window position and size
- ✅ **Wails Integration** - Clean API between Go and React

## Database Location

The database is stored in the user's config directory:

- **macOS**: `~/Library/Application Support/V2_byteframes/app.db`
- **Linux**: `~/.config/V2_byteframes/app.db`
- **Windows**: `%APPDATA%\V2_byteframes\app.db`

## Schema

### Tables

#### `settings`
Stores key-value configuration settings.

```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

**Usage:**
- Theme preference: `key="theme"`, `value="dockview-theme-dark"`
- Any app-wide settings

#### `layouts`
Stores saved Dockview layout configurations.

```sql
CREATE TABLE layouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    layout_json TEXT NOT NULL,
    is_active INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

**Usage:**
- Save current panel arrangement
- Load previous layouts
- One layout marked as active (auto-loaded on startup)

#### `window_state`
Stores window position and size for next launch.

```sql
CREATE TABLE window_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    x INTEGER,
    y INTEGER,
    width INTEGER,
    height INTEGER,
    maximized INTEGER DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

**Note:** Currently not implemented in the UI, but database schema is ready.

#### `schema_version`
Tracks current database version for migrations.

```sql
CREATE TABLE schema_version (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    version INTEGER NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

## Go API

### Database Package (`db/`)

#### Initialize Database
```go
import "V2_byteframes/db"

database, err := db.New(appDataDir)
if err != nil {
    log.Fatal(err)
}
defer database.Close()
```

#### Settings Methods
```go
// Get a setting
value, err := database.GetSetting("theme")

// Set a setting
err := database.SetSetting("theme", "dockview-theme-light")

// Get all settings
settings, err := database.GetAllSettings()
```

#### Layout Methods
```go
// Save a layout
layout, err := database.SaveLayout("My Layout", layoutJSON)

// Get all layouts
layouts, err := database.GetAllLayouts()

// Get specific layout
layout, err := database.GetLayout(id)

// Set active layout
err := database.SetActiveLayout(id)

// Get active layout
layout, err := database.GetActiveLayout()

// Delete layout
err := database.DeleteLayout(id)
```

#### Window State Methods
```go
// Save window state
err := database.SaveWindowState(x, y, width, height, maximized)

// Get window state
state, err := database.GetWindowState()
```

### Wails Service Layer (`app.go`)

The `App` struct exposes database methods to the frontend:

```go
type App struct {
    ctx context.Context
    db  *db.Database
}

// These methods are automatically bound to the frontend
func (a *App) GetSetting(key string) (string, error)
func (a *App) SetSetting(key, value string) error
func (a *App) SaveLayout(name, layoutJSON string) (*models.Layout, error)
func (a *App) GetAllLayouts() ([]models.Layout, error)
// ... and more
```

## React/TypeScript API

Wails automatically generates TypeScript bindings in `frontend/wailsjs/go/main/App.js`.

### Import
```typescript
import { 
    GetSetting, 
    SetSetting, 
    SaveLayout, 
    GetAllLayouts, 
    GetActiveLayout,
    SetActiveLayout 
} from "../wailsjs/go/main/App";
```

### Usage Examples

#### Save Theme Preference
```typescript
import { SetSetting } from "../wailsjs/go/main/App";

const saveTheme = async (themeName: string) => {
    await SetSetting("theme", themeName);
};
```

#### Load Theme Preference
```typescript
import { GetSetting } from "../wailsjs/go/main/App";

const loadTheme = async () => {
    const savedTheme = await GetSetting("theme");
    if (savedTheme) {
        // Apply theme
    }
};
```

#### Save Current Layout
```typescript
import { SaveLayout } from "../wailsjs/go/main/App";

const saveCurrentLayout = async (api: DockviewApi, name: string) => {
    const layoutData = api.toJSON();
    const layoutJSON = JSON.stringify(layoutData);
    await SaveLayout(name, layoutJSON);
};
```

#### Load Saved Layouts
```typescript
import { GetAllLayouts, SetActiveLayout } from "../wailsjs/go/main/App";

const loadLayouts = async () => {
    const layouts = await GetAllLayouts();
    // Display in UI
};

const applyLayout = async (api: DockviewApi, layout: Layout) => {
    const layoutData = JSON.parse(layout.layout_json);
    api.fromJSON(layoutData);
    await SetActiveLayout(layout.id);
};
```

#### Load Active Layout on Startup
```typescript
import { GetActiveLayout } from "../wailsjs/go/main/App";

const onReady = async (event: DockviewReadyEvent) => {
    const activeLayout = await GetActiveLayout();
    if (activeLayout && activeLayout.layout_json) {
        const layoutData = JSON.parse(activeLayout.layout_json);
        event.api.fromJSON(layoutData);
    }
};
```

## Implementation Details

### How Theme Persistence Works

1. **On Theme Change**:
   ```typescript
   useEffect(() => {
       if (currentTheme) {
           SetSetting("theme", currentTheme.name);
       }
   }, [currentTheme]);
   ```

2. **On App Load**:
   ```typescript
   useEffect(() => {
       GetSetting("theme").then((savedTheme) => {
           if (savedTheme) {
               const theme = themes.find((t) => t.theme.name === savedTheme);
               if (theme) setCurrentTheme(theme.theme);
           }
       });
   }, []);
   ```

### How Layout Persistence Works

1. **Save Layout**:
   - User clicks "💾 Save"
   - Enters layout name in modal
   - Current Dockview layout serialized with `api.toJSON()`
   - Saved to database via `SaveLayout(name, json)`

2. **Load Layout**:
   - User clicks "📂 Load"
   - Displays all saved layouts from database
   - Click on layout to apply
   - Deserialize with `api.fromJSON(data)`
   - Mark as active with `SetActiveLayout(id)`

3. **Auto-load on Startup**:
   - `onReady` handler calls `GetActiveLayout()`
   - If found, deserialize and apply
   - Otherwise, create default panels

## Adding New Tables

1. **Create migration** in `db/migrations.go`:
   ```go
   func (d *Database) migration002_myFeature() error {
       _, err := d.DB.Exec(`
           CREATE TABLE my_table (
               id INTEGER PRIMARY KEY,
               data TEXT NOT NULL
           )
       `)
       return err
   }
   ```

2. **Add to migrations list**:
   ```go
   migrations := []migration{
       {version: 1, up: d.migration001_initial},
       {version: 2, up: d.migration002_myFeature}, // Add here
   }
   ```

3. **Create model** in `models/models.go`:
   ```go
   type MyData struct {
       ID   int    `json:"id"`
       Data string `json:"data"`
   }
   ```

4. **Add queries** in `db/queries.go`:
   ```go
   func (d *Database) GetMyData(id int) (*models.MyData, error) {
       // Implementation
   }
   ```

5. **Expose to frontend** in `app.go`:
   ```go
   func (a *App) GetMyData(id int) (*models.MyData, error) {
       return a.db.GetMyData(id)
   }
   ```

6. **Regenerate bindings**:
   ```bash
   wails generate module
   ```

## Migration System

Migrations run automatically on app startup in `main.go`:

```go
database, err := db.New(appDataDir)
// Migrations execute here
```

The system:
- Checks current version in `schema_version` table
- Runs migrations newer than current version
- Updates version after each successful migration
- Transactions ensure atomic updates

## Error Handling

All database operations return errors. Handle them appropriately:

```typescript
try {
    await SaveLayout(name, json);
    console.log("Layout saved successfully");
} catch (err) {
    console.error("Failed to save layout:", err);
    // Show error to user
}
```

## Best Practices

1. **Always check for errors** - Database operations can fail
2. **Use transactions** for multi-step operations (already done in migrations)
3. **Validate input** before saving to database
4. **Keep JSON serialization simple** - Don't store complex objects
5. **Test migrations** with existing databases before deploying

## Performance Notes

- SQLite configured with `SetMaxOpenConns(1)` - optimal for single-user desktop app
- Indexes created on frequently queried columns (`layouts.is_active`)
- Migrations run only once per version
- Pure Go implementation means no CGO overhead

## Debugging

To inspect the database:

```bash
# macOS
sqlite3 ~/Library/Application\ Support/V2_byteframes/app.db

# Common queries
.tables                          # List all tables
SELECT * FROM settings;          # View settings
SELECT * FROM layouts;           # View layouts
SELECT version FROM schema_version; # Check version
```

## Future Enhancements

Possible additions:
- Window state save/restore (schema ready, UI not implemented)
- Layout export/import (JSON files)
- Layout snapshots/versioning
- User profiles with separate settings
- Recent layouts list
- Layout thumbnails/previews

---

This database integration provides a solid foundation for any data persistence needs in your Wails application.
