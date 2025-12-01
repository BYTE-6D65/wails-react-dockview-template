import React, { useState } from "react";
import {
  DockviewReact,
  DockviewReadyEvent,
  IDockviewPanelProps,
  DockviewApi,
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
import "dockview-core/dist/styles/dockview.css";
import { TestPanel } from "./TestPanel";
import "./panels.css";

const components = {
  default: (props: IDockviewPanelProps) => {
    return <TestPanel {...props} />;
  },
};

const addDefaultPanels = (api: DockviewApi) => {
  api.addPanel({
    id: "panel_1",
    component: "default",
    title: "Panel 1",
  });
  api.addPanel({
    id: "panel_2",
    component: "default",
    title: "Panel 2",
    position: { referencePanel: "panel_1", direction: "right" },
  });
};

export const DockviewLayout: React.FC = () => {
  const themes: Array<{ name: string; theme: DockviewTheme }> = [
    { name: "Dark", theme: themeDark },
    { name: "Light", theme: themeLight },
    { name: "Visual Studio", theme: themeVisualStudio },
    { name: "Abyss", theme: themeAbyss },
    { name: "Dracula", theme: themeDracula },
    { name: "Replit", theme: themeReplit },
    { name: "Abyss Spaced", theme: themeAbyssSpaced },
    { name: "Light Spaced", theme: themeLightSpaced },
  ];

  const [currentTheme, setCurrentTheme] = useState<DockviewTheme>(themeDark);
  const [showThemeMenu, setShowThemeMenu] = useState(false);
  const [api, setApi] = useState<DockviewApi | null>(null);
  const [panelCounter, setPanelCounter] = useState(3);

  const onReady = (event: DockviewReadyEvent) => {
    setApi(event.api);
    addDefaultPanels(event.api);
  };

  const addNewPanel = () => {
    if (!api) return;

    const newPanelId = `panel_${panelCounter}`;
    api.addPanel({
      id: newPanelId,
      component: "default",
      title: `Panel ${panelCounter}`,
    });
    setPanelCounter(panelCounter + 1);
  };

  return (
    <div style={{ width: "100%", height: "100%", position: "relative" }}>
      <div
        style={{
          position: "absolute",
          top: "10px",
          right: "10px",
          zIndex: 10000,
          display: "flex",
          gap: "8px",
        }}
      >
        <button
          onClick={addNewPanel}
          style={{
            background: "#28a745",
            color: "white",
            border: "none",
            borderRadius: "4px",
            padding: "8px 12px",
            cursor: "pointer",
          }}
        >
          ➕ Add Panel
        </button>
        <button
          onClick={() => setShowThemeMenu(!showThemeMenu)}
          style={{
            background: "#6f42c1",
            color: "white",
            border: "none",
            borderRadius: "4px",
            padding: "8px 12px",
            cursor: "pointer",
          }}
        >
          🎨 Theme
        </button>
      </div>

      {showThemeMenu && (
        <div
          style={{
            position: "absolute",
            top: "50px",
            right: "10px",
            zIndex: 10000,
            background: "#1e1e1e",
            border: "1px solid #333",
            borderRadius: "4px",
            padding: "8px",
            minWidth: "200px",
          }}
        >
          {themes.map((themeOption) => (
            <div
              key={themeOption.name}
              onClick={() => {
                setCurrentTheme(themeOption.theme);
                setShowThemeMenu(false);
              }}
              style={{
                padding: "8px 12px",
                cursor: "pointer",
                color: "white",
                background:
                  currentTheme.name === themeOption.theme.name
                    ? "#6f42c1"
                    : "transparent",
              }}
            >
              {themeOption.name}
            </div>
          ))}
        </div>
      )}

      <div style={{ width: "100%", height: "100%" }}>
        <DockviewReact
          components={components}
          onReady={onReady}
          theme={currentTheme}
        />
      </div>
    </div>
  );
};
