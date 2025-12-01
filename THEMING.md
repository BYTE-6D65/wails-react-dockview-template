# Dockview Theming Guide

## Overview

The V2_byteframes project uses Dockview's built-in theming system to provide consistent visual styling across all panels and layout components.

## How Theming Works

### 1. Required CSS Import

The core Dockview stylesheet must be imported before any components can be themed:

```typescript
import "dockview-core/dist/styles/dockview.css";
```

This CSS file contains all the CSS variables and theme classes that power the theming system.

### 2. Built-in Theme Objects

Dockview provides pre-configured theme objects that you can import and use:

```typescript
import {
  DockviewTheme,
  themeDark,
  themeLight,
  themeAbyss,
  themeDracula,
  themeReplit,
  themeVisualStudio,
  themeAbyssSpaced,
  themeLightSpaced,
} from "dockview-react";
```

Each theme object has this structure:

```typescript
interface DockviewTheme {
  name: string;           // Theme identifier (e.g., "dockview-theme-light")
  className: string;      // CSS class name to apply
  gap?: number;          // Optional spacing between groups
  dndOverlayMounting?: 'relative' | 'absolute';  // Drag overlay positioning
  dndPanelOverlay?: 'content' | 'group';         // Drag overlay scope
}
```

### 3. Applying Themes to DockviewReact

For `DockviewReact` components, use the `theme` prop (NOT `className`):

```typescript
<DockviewReact
  components={components}
  onReady={onReady}
  theme={currentTheme}  // Pass theme object here
/>
```

**Important:** The `className` prop is only for other Dockview components like Splitview, Gridview, Paneview, and Tabview.

### 4. Theme Inheritance

When you set a theme on the DockviewReact component:
- All panels automatically inherit the theme
- New panels added dynamically will use the current theme
- No manual styling is needed for individual panels

The theme affects:
- Tab bar background and text colors
- Panel backgrounds
- Border colors
- Hover states
- Active/inactive states
- Separator colors

## Implementation in V2_byteframes

### Current Setup

Located in `V2_byteframes/frontend/src/DockviewLayout.tsx`:

```typescript
import {
  DockviewReact,
  DockviewTheme,
  themeDark,
  themeLight,
} from "dockview-react";
import "dockview-core/dist/styles/dockview.css";

export const DockviewLayout: React.FC = () => {
  const themes = [
    { name: "Dark", theme: themeDark },
    { name: "Light", theme: themeLight },
  ];

  const [currentTheme, setCurrentTheme] = useState<DockviewTheme>(themeDark);

  return (
    <DockviewReact
      components={components}
      onReady={onReady}
      theme={currentTheme}
    />
  );
};
```

### Theme Switcher

The theme switcher UI allows runtime theme changes:

1. Click the "🎨 Theme" button in the top-right corner
2. Select "Dark" or "Light" from the dropdown
3. The theme updates immediately for all panels

## Using CSS Variables in Custom Panels

Panels can use Dockview's CSS variables to match the current theme:

```css
.bf-panel {
  color: var(--dv-activegroup-visiblepanel-tab-color);
  background-color: var(--dv-group-view-background-color);
}
```

### Available CSS Variables

Common Dockview CSS variables include:

- `--dv-group-view-background-color` - Panel background
- `--dv-tabs-and-actions-container-background-color` - Tab bar background
- `--dv-activegroup-visiblepanel-tab-color` - Active tab text color
- `--dv-active-sash-color` - Active separator color
- `--dv-separator-border` - Border between panels
- `--dv-paneview-header-border-color` - Header borders

These variables automatically update when the theme changes.

## Adding More Themes

To add additional built-in themes:

```typescript
import {
  themeAbyss,
  themeDracula,
  themeReplit,
} from "dockview-react";

const themes = [
  { name: "Dark", theme: themeDark },
  { name: "Light", theme: themeLight },
  { name: "Abyss", theme: themeAbyss },
  { name: "Dracula", theme: themeDracula },
  { name: "Replit", theme: themeReplit },
];
```

## Creating Custom Themes

You can create custom themes by defining a DockviewTheme object:

```typescript
const myCustomTheme: DockviewTheme = {
  name: "my-custom-theme",
  className: "my-custom-theme-class",
  gap: 5,  // 5px gap between groups
};
```

Then define the CSS class with your custom variables:

```css
.my-custom-theme-class {
  --dv-group-view-background-color: #1a1a2e;
  --dv-tabs-and-actions-container-background-color: #16213e;
  --dv-activegroup-visiblepanel-tab-color: #eef4ed;
  /* ... other variables ... */
}
```

## Troubleshooting

### Themes Not Working

1. **Missing CSS import** - Ensure `import "dockview-core/dist/styles/dockview.css";` is present
2. **Wrong prop** - Use `theme={themeObject}` NOT `className="theme-class"`
3. **Case sensitivity** - The div ID must match CSS selector (`#App` needs `#App` in CSS)

### Panels Not Visible

1. **Container height** - Ensure parent div has explicit height (e.g., `height: 100vh`)
2. **CSS file path** - Verify the dockview.css path is correct for your node_modules structure

## References

- [Dockview Documentation](https://dockview.dev/docs/overview/getStarted/theme)
- [Dockview GitHub](https://github.com/mathuo/dockview)
- [Theme SCSS Source](https://github.com/mathuo/dockview/blob/master/packages/dockview-core/src/theme.scss)
