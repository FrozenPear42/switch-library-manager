import { Switch } from "react-aria-components";
import styles from "./Toggle.module.css";

type ToggleProps = {
  label?: string;
};

export default function Toggle({ label }: ToggleProps) {
  return (
    <Switch className={styles.switch}>
      <div className={styles.indicator} />
      {label}
    </Switch>
  );
}
