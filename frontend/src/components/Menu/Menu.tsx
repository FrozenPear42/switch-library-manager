import {
  IconDatabase,
  IconLibrary,
  IconNut,
  IconSettings,
} from "@tabler/icons-react";
import styles from "./Menu.module.css";
import classNames from "classnames";

export default function Menu() {
  return (
    <div className={styles.menu}>
      <a className={classNames(styles.menuItem, styles.active)}>
        <IconLibrary />
        <label>Library</label>
      </a>
      <a className={styles.menuItem}>
        <IconDatabase />
        <label>Catalogue</label>
      </a>
      <a className={styles.menuItem}>
        <IconNut />
        <label>NUT</label>
      </a>
      <a className={styles.menuItem}>
        <IconSettings />
        <label>Settings</label>
      </a>
    </div>
  );
}
