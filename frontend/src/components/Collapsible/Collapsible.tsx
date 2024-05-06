import classNames from "classnames";
import styles from "./Collapsible.module.css";
import { useRef, useState } from "react";

type CollapsibleProps = React.PropsWithChildren<{
  startOpened: boolean;
  openText?: string;
  closeText?: string;
}>;

export default function Collapsible({
  startOpened,
  children,
  openText = "Show",
  closeText = "Hide",
}: CollapsibleProps) {
  const [opened, setOpened] = useState<boolean>(startOpened);
  const contentRef = useRef<HTMLDivElement>(null);

  const onToggle = () => {
    // console.log(
    //   contentRef.current?.getBoundingClientRect(),
    //   contentRef.current?.offsetHeight,
    //   opened
    // );
    // if (contentRef.current && !opened) {
    //   contentRef.current.style.setProperty(
    //     "--content-height",
    //     `${contentRef.current.getBoundingClientRect().height}px`
    //   );
    // }
    setOpened(!opened);
  };

  return (
    <div className={styles.wrapper}>
      <div
        className={classNames(styles.content, !opened && styles.collapsed)}
        ref={contentRef}
      >
        {children}
      </div>
      <div className={classNames(styles.action, !opened && styles.collapsed)}>
        <button type="button" onClick={onToggle} className={styles.button}>
          {!opened ? openText : closeText}
        </button>
      </div>
    </div>
  );
}
