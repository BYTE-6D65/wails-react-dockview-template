import React, { useState, useEffect } from "react";
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
import {
  GetSetting,
  SetSetting,
  SaveLayout,
  GetAllLayouts,
  GetActiveLayout,
  SetActiveLayout,
} from "../wailsjs/go/main/App";

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

interface Layout {
  id: number;
  name: string;
  layout_json: string;
  is_active: boolean;
  created_at: number;
  updated_at: number;
}

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
  const [showSaveModal, setShowSaveModal] = useState(false);
  const [showLoadModal, setShowLoadModal] = useState(false);
  const [layoutName, setLayoutName] = useState("");
  const [savedLayouts, setSavedLayouts] = useState<Layout[]>([]);

  // Load saved theme on mount
  useEffect(() => {
    GetSetting("theme").then((savedTheme) => {
      if (savedTheme) {
        const theme = themes.find((t) => t.theme.name === savedTheme);
        if (theme) {
          setCurrentTheme(theme.theme);
        }
      }
    });
  }, []);

  // Save theme when it changes
  useEffect(() => {
    if (currentTheme) {
      SetSetting("theme", currentTheme.name);
    }
  }, [currentTheme]);

  const onReady = async (event: DockviewReadyEvent) => {
    setApi(event.api);

    // Try to load active layout
    try {
      const activeLayout = await GetActiveLayout();
      if (activeLayout && activeLayout.layout_json) {
        const layoutData = JSON.parse(activeLayout.layout_json);
        event.api.fromJSON(layoutData);
        return;
      }
    } catch (err) {
      console.log("No active layout found, using defaults");
    }

    // Fall back to default panels
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

  const handleSaveLayout = async () => {
    if (!api || !layoutName.trim()) return;

    try {
      const layoutData = api.toJSON();
      const layoutJSON = JSON.stringify(layoutData);
      await SaveLayout(layoutName, layoutJSON);
      setShowSaveModal(false);
      setLayoutName("");
      console.log("Layout saved:", layoutName);
    } catch (err) {
      console.error("Failed to save layout:", err);
    }
  };

  const handleLoadLayoutClick = async () => {
    try {
      const layouts = await GetAllLayouts();
      setSavedLayouts(layouts || []);
      setShowLoadModal(true);
    } catch (err) {
      console.error("Failed to load layouts:", err);
    }
  };

  const handleLoadLayout = async (layout: Layout) => {
    if (!api) return;

    try {
      const layoutData = JSON.parse(layout.layout_json);
      api.fromJSON(layoutData);
      await SetActiveLayout(layout.id);
      setShowLoadModal(false);
      console.log("Layout loaded:", layout.name);
    } catch (err) {
      console.error("Failed to load layout:", err);
    }
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
          onClick={handleLoadLayoutClick}
          style={{
            background: "#007bff",
            color: "white",
            border: "none",
            borderRadius: "4px",
            padding: "8px 12px",
            cursor: "pointer",
          }}
        >
          📂 Load
        </button>
        <button
          onClick={() => setShowSaveModal(true)}
          style={{
            background: "#17a2b8",
            color: "white",
            border: "none",
            borderRadius: "4px",
            padding: "8px 12px",
            cursor: "pointer",
          }}
        >
          💾 Save
        </button>
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

      {/* Save Layout Modal */}
      {showSaveModal && (
        <div
          style={{
            position: "fixed",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: "rgba(0, 0, 0, 0.7)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 20000,
          }}
          onClick={() => setShowSaveModal(false)}
        >
          <div
            style={{
              background: "#1e1e1e",
              padding: "24px",
              borderRadius: "8px",
              border: "1px solid #333",
              minWidth: "400px",
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <h3 style={{ margin: "0 0 16px 0", color: "white" }}>
              Save Layout
            </h3>
            <input
              type="text"
              value={layoutName}
              onChange={(e) => setLayoutName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleSaveLayout();
                if (e.key === "Escape") setShowSaveModal(false);
              }}
              placeholder="Layout name"
              autoFocus
              style={{
                width: "100%",
                padding: "10px",
                background: "#252526",
                border: "1px solid #333",
                borderRadius: "4px",
                color: "white",
                fontSize: "14px",
                marginBottom: "16px",
              }}
            />
            <div
              style={{
                display: "flex",
                gap: "8px",
                justifyContent: "flex-end",
              }}
            >
              <button
                onClick={() => setShowSaveModal(false)}
                style={{
                  padding: "8px 16px",
                  background: "#6c757d",
                  color: "white",
                  border: "none",
                  borderRadius: "4px",
                  cursor: "pointer",
                }}
              >
                Cancel
              </button>
              <button
                onClick={handleSaveLayout}
                disabled={!layoutName.trim()}
                style={{
                  padding: "8px 16px",
                  background: layoutName.trim() ? "#17a2b8" : "#444",
                  color: "white",
                  border: "none",
                  borderRadius: "4px",
                  cursor: layoutName.trim() ? "pointer" : "not-allowed",
                }}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Load Layout Modal */}
      {showLoadModal && (
        <div
          style={{
            position: "fixed",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: "rgba(0, 0, 0, 0.7)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 20000,
          }}
          onClick={() => setShowLoadModal(false)}
        >
          <div
            style={{
              background: "#1e1e1e",
              padding: "24px",
              borderRadius: "8px",
              border: "1px solid #333",
              minWidth: "400px",
              maxHeight: "80vh",
              overflow: "auto",
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <h3 style={{ margin: "0 0 16px 0", color: "white" }}>
              Load Layout
            </h3>
            {savedLayouts.length === 0 ? (
              <p style={{ color: "#888", marginBottom: "16px" }}>
                No saved layouts found. Save your current layout first!
              </p>
            ) : (
              <div
                style={{
                  display: "flex",
                  flexDirection: "column",
                  gap: "8px",
                  marginBottom: "16px",
                }}
              >
                {savedLayouts.map((layout) => (
                  <div
                    key={layout.id}
                    onClick={() => handleLoadLayout(layout)}
                    style={{
                      padding: "12px",
                      background: layout.is_active ? "#1a4d2e" : "#252526",
                      border: layout.is_active
                        ? "1px solid #28a745"
                        : "1px solid #333",
                      borderRadius: "4px",
                      cursor: "pointer",
                      transition: "all 0.2s",
                    }}
                    onMouseEnter={(e) => {
                      if (!layout.is_active) {
                        e.currentTarget.style.background = "#2a2d2e";
                        e.currentTarget.style.borderColor = "#555";
                      }
                    }}
                    onMouseLeave={(e) => {
                      if (!layout.is_active) {
                        e.currentTarget.style.background = "#252526";
                        e.currentTarget.style.borderColor = "#333";
                      }
                    }}
                  >
                    <div
                      style={{
                        color: "white",
                        fontSize: "14px",
                        fontWeight: 600,
                        marginBottom: "4px",
                      }}
                    >
                      {layout.name}
                      {layout.is_active && (
                        <span
                          style={{
                            marginLeft: "8px",
                            color: "#28a745",
                            fontSize: "12px",
                          }}
                        >
                          ● Active
                        </span>
                      )}
                    </div>
                    <div style={{ color: "#888", fontSize: "12px" }}>
                      {new Date(layout.updated_at * 1000).toLocaleString()}
                    </div>
                  </div>
                ))}
              </div>
            )}
            <button
              onClick={() => setShowLoadModal(false)}
              style={{
                padding: "8px 16px",
                background: "#6c757d",
                color: "white",
                border: "none",
                borderRadius: "4px",
                cursor: "pointer",
              }}
            >
              Close
            </button>
          </div>
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
