import {
  IconAlertSquareRoundedFilled,
  IconAlertTriangleFilled,
  IconExclamationMark,
} from "@tabler/icons-react";
import styles from "./Header.module.css";
import logo from "../../assets/images/logo.png";

export default function Header() {
  return (
    <div className={styles.header}>
      <div className={styles.brand}>
        <img src={logo} className={styles.brandIcon} />
        <div className="d-grid">
          <div className={styles.brandName}>Switch Library Manager</div>
          <div className={styles.brandVersion}>v1.2</div>
        </div>
      </div>
      <div className={styles.stats}>
        <div>Titles: 10, DLCs: 120</div>
        <div>
          <IconAlertSquareRoundedFilled /> 8 issues
        </div>
        {/* <div className={styles.actions}>
          <button>Reload</button>
          <button>Hard reload</button>
        </div> */}
      </div>
    </div>
  );
}
