import { ProgressBar, ProgressBarProps, Label } from "react-aria-components";

import styles from "./Progress.module.css";

type ProgressProps = ProgressBarProps & {
  label?: string;
  details?: string;
};

export default function Progress({ label, details, ...props }: ProgressProps) {
  return (
    <ProgressBar className={styles.progressBar} {...props}>
      {({ percentage, valueText }) => (
        <>
          <Label>{label}</Label>
          <div className={styles.value}>{valueText}</div>
          <div className={styles.bar}>
            <div
              className={styles.fill}
              style={{ transform: `translateX(-${100 - (percentage || 0)}%)` }}
            ></div>
          </div>
          <div className={styles.details}>{details}</div>
        </>
      )}
    </ProgressBar>
  );
}
