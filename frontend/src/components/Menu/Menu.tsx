import {
  IconDatabase,
  IconFile,
  IconLibrary,
  IconNut,
  IconSettings,
} from "@tabler/icons-react";
import styles from "./Menu.module.css";
import classNames from "classnames";
import { Link } from "wouter";

export default function Menu() {
  return (
    <div className={styles.menu}>
      <Link
        to="/library"
        className={(isActive) =>
          classNames(styles.menuItem, isActive && styles.active)
        }
      >
        <IconLibrary />
        <span>Library</span>
      </Link>

      <Link
        to="/files"
        className={(isActive) =>
          classNames(styles.menuItem, isActive && styles.active)
        }
      >
        <IconFile />
        <span>Files</span>
      </Link>

      <Link
        to="/catalog"
        className={(isActive) =>
          classNames(styles.menuItem, isActive && styles.active)
        }
      >
        <IconDatabase />
        <span>Catalog</span>
      </Link>

      <Link
        to="/nut"
        className={(isActive) =>
          classNames(styles.menuItem, isActive && styles.active)
        }
      >
        <IconNut />
        <span>NUT</span>
      </Link>

      <Link
        to="/settings"
        className={(isActive) =>
          classNames(styles.menuItem, isActive && styles.active)
        }
      >
        <IconSettings />
        <span>Settings</span>
      </Link>
    </div>
  );
}
