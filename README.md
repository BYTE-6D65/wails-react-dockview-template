# Wails + React + TypeScript + Dockview Template

A production-ready template for building desktop applications with **Wails**, **React 19**, **TypeScript 5.9**, and **Dockview** theming system.

This template demonstrates the **correct** way to implement Dockview theming with all 8 built-in themes working perfectly.

## 🎯 Features

- ✅ **Wails v2** - Build desktop apps with Go backend
- ✅ **React 19.2** - Latest stable React
- ✅ **TypeScript 5.9** - Full type safety
- ✅ **Vite 7.2** - Lightning-fast HMR
- ✅ **SQLite Database** - Pure Go (modernc.org/sqlite, no CGO)
- ✅ **Dockview 4.11** - Advanced layout system with proper theming
- ✅ **8 Built-in Themes** - Dark, Light, Visual Studio, Abyss, Dracula, Replit, and Spaced variants
- ✅ **Theme Persistence** - Remembers your theme choice
- ✅ **Layout Save/Load** - Save and restore panel arrangements
- ✅ **Theme Hot-Switching** - Change themes at runtime
- ✅ **Dynamic Panel Management** - Add/remove panels on the fly
- ✅ **Automatic Theme Inheritance** - New panels automatically match current theme

## 📦 Tech Stack

| Package | Version |
|---------|---------|
| React | 19.2.0 |
| React DOM | 19.2.0 |
| TypeScript | 5.9.3 |
| Vite | 7.2.4 |
| Dockview | 4.11.0 |
| Dockview React | 4.11.0 |
| SQLite (modernc.org/sqlite) | 1.40.1 |
| @vitejs/plugin-react | 5.1.1 |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Wails CLI v2

### Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd V2_byteframes

# Install frontend dependencies
cd V2_byteframes/frontend
npm install

# Return to root
cd ..

# Run in development mode
wails dev
```

### Build for Production

```bash
wails build
```

## 🎨 Theming Implementation

This template demonstrates the **correct** way to implement Dockview theming. See [THEMING.md](./THEMING.md) for detailed documentation.

### Key Principles

1. **Import the core CSS**
```typescript
import "dockview-core/dist/styles/dockview.css";
```

2. **Use built-in theme objects**
```typescript
import { themeDark, themeLight, themeAbyss } from "dockview-react";
```

3. **Pass theme to DockviewReact**
```typescript
<DockviewReact theme={currentTheme} />
```

4. **Use Dockview CSS variables in custom components**
```css
.my-panel {
  color: var(--dv-activegroup-visiblepanel-tab-color);
  background-color: var(--dv-group-view-background-color);
}
```

## 📚 Documentation

- [THEMING.md](./THEMING.md) - Complete theming guide
- [DATABASE.md](./DATABASE.md) - SQLite integration and API reference
- [LESSONS_LEARNED.md](./LESSONS_LEARNED.md) - Common mistakes and how to avoid them

## 🏗️ Project Structure

```
V2_byteframes/
├── V2_byteframes/
│   ├── frontend/
│   │   ├── src/
│   │   │   ├── App.tsx              # Main app component
│   │   │   ├── DockviewLayout.tsx   # Dockview configuration
│   │   │   ├── TestPanel.tsx        # Example panel component
│   │   │   ├── App.css              # App styles
│   │   │   ├── panels.css           # Panel styles with Dockview variables
│   │   │   └── main.tsx             # React entry point
│   │   ├── package.json
│   │   └── vite.config.ts
│   ├── main.go                      # Go backend
│   └── wails.json                   # Wails configuration
├── THEMING.md                       # Theming documentation
├── LESSONS_LEARNED.md               # V1 mistakes analysis
└── README.md                        # This file
```

## 🎯 Available Themes

1. **Dark** - Basic dark theme
2. **Light** - Basic light theme  
3. **Visual Studio** - VS-inspired theme with blue accents
4. **Abyss** - Deep blue palette with purple accents
5. **Dracula** - Popular dark theme
6. **Replit** - Light design with rounded corners
7. **Abyss Spaced** - Abyss with extra panel spacing
8. **Light Spaced** - Light with extra panel spacing

## 🧪 Testing the Template

1. **Start the app**: `wails dev`
2. **Test theme switching**: Click "🎨 Theme" button and try all 8 themes
3. **Test theme persistence**: Change theme, restart app, verify theme is remembered
4. **Test layout save**: Arrange panels, click "💾 Save", enter name
5. **Test layout load**: Click "📂 Load", select saved layout
6. **Test dynamic panels**: Click "➕ Add Panel" to add new panels
7. **Verify theme inheritance**: Add panels while on different themes

All themes should apply instantly to all panels, including newly created ones. Theme choice and active layout persist across app restarts.

## ❌ Common Mistakes to Avoid

### DON'T ❌
```typescript
// Don't apply className to wrapper div
<div className="dockview-theme-dark">
  <DockviewReact />
</div>

// Don't use className prop on DockviewReact
<DockviewReact className="dockview-theme-dark" />

// Don't hardcode colors
background: #1e1e1e;

// Don't use non-existent variables
var(--dv-text-color)
```

### DO ✅
```typescript
// Pass theme object to DockviewReact
<DockviewReact theme={themeDark} />

// Use real Dockview CSS variables
background: var(--dv-group-view-background-color);
```

See [LESSONS_LEARNED.md](./LESSONS_LEARNED.md) for detailed analysis of what went wrong in previous implementations.

## 🤝 Contributing

This is a template repository. Feel free to:
- Fork it for your own projects
- Submit issues for bugs or improvements
- Create PRs with enhancements

## 📄 License

MIT License - see LICENSE file for details

## 🔗 Resources

- [Wails Documentation](https://wails.io)
- [Dockview Documentation](https://dockview.dev)
- [Dockview Theming Guide](https://dockview.dev/docs/overview/getStarted/theme)
- [React 19 Documentation](https://react.dev)
- [Vite Documentation](https://vitejs.dev)

## 🙏 Acknowledgments

- **Wails** - For making desktop app development with Go/React enjoyable
- **Dockview** - For the excellent layout system
- Built with insights from debugging V1 theming issues (see LESSONS_LEARNED.md)

---

**Note**: This template represents the **correct** implementation of Dockview theming after learning from V1 mistakes. All 8 themes work perfectly, and theme inheritance is automatic.
