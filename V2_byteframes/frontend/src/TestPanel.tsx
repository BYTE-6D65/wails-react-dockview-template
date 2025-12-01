import { IDockviewPanelProps } from "dockview-react";
import "./panels.css";

export const TestPanel = (props: IDockviewPanelProps) => {
  return (
    <div className="bf-panel">
      <div className="bf-panel-content">
        <h3>{props.api.title}</h3>
        <p>This is a test panel. The content should be themed correctly.</p>
      </div>
    </div>
  );
};
