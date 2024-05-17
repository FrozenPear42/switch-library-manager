import { StartupState } from "../../hooks/useStartup";
import Progress from "../Progress/Progress";
import styles from "./LoadingPanel.module.css";
import logo from "../../assets/images/logo.png";

type LoadingPanelProps = {
  state: StartupState;
};

export default function LoadingPanel({ state }: LoadingPanelProps) {
  return (
    <div className={styles.panel}>
      <div className={styles.logo}>
        <img src={logo} className={styles.logoImage} />
        <div className={styles.logoName}>Switch Library Manager</div>
        <div className={styles.logoVersion}>v1.4</div>
      </div>

      <div className={styles.progressBar}>
        <Progress
          label={"Loading..."}
          value={state ? (state.stageCurrent / state.stageTotal) * 100 : 0}
          details={state?.stageMessage || "\u2000"}
        ></Progress>
      </div>
    </div>
  );
}
