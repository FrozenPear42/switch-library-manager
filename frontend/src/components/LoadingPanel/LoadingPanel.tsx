import { useStartup } from "../../hooks/useStartup";
import ProgressBar from "../ProgressBar/ProgressBar";
import styles from "./LoadingPanel.module.css";

export default function LoadingPanel() {
  const { state } = useStartup();

  return (
    <div className={styles.panel}>
      <div className={styles.logo}>
        <img
          src={"https://placehold.co/400x400/orange/white"}
          className={styles.logoImage}
        />
        <div className={styles.logoName}>Switch Library Manager</div>
        <div className={styles.logoVersion}>v1.4</div>
      </div>

      <div className={styles.progressBar}>
        <ProgressBar
          progress={state ? ((state.stageCurrent / state.stageTotal) * 100) : 0}
        ></ProgressBar>
      </div>
      <div className={styles.progressText}>
        {state?.stageMessage || "\u2000"}
      </div>
    </div>
  );
}
