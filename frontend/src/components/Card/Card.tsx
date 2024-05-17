import classNames from "classnames";
import { HTMLAttributes, forwardRef } from "react";

import styles from "./Card.module.css";

type CardProps = HTMLAttributes<HTMLDivElement> & {
  coverImage?: string;
};

export const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, coverImage, children, ...props }, ref) => (
    <div className={classNames(styles.card, className)} {...props} ref={ref}>
      {coverImage && <img src={coverImage} className={styles.coverImage} />}
      {children}
    </div>
  )
);
