import {
  IconExclamationCircle,
  IconMessageCircleExclamation,
} from "@tabler/icons-react";
import Spinner from "../../components/Spinner/Spinner";
import { useLibrary } from "../../hooks/useLibrary";
import GameCard from "./GameCard";
import styles from "./Library.module.css";
import { AppSelect, AppSelectItem } from "../../components/Select/Select";

export default function Library() {
  const { data: games, isLoading, error } = useLibrary();

  if (isLoading) {
    return (
      <div className={styles.page}>
        <div className={styles.loading}>
          <Spinner />
          <div>Loading library...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.page}>
        <div className={styles.loadingError}>
          <IconExclamationCircle size={60} stroke={2} />
          <div>Could not load</div>
          <div>{`${error}`}</div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <div className={styles.actions}>
        <AppSelect
          label="Sort by"
          items={[
            { key: "id", label: "ID" },
            { key: "name", label: "Name" },
            { key: "region", label: "Region" },
            { key: "issues", label: "Issues" },
          ]}
          defaultSelectedKey={"name"}
        >
          {(item) => (
            <AppSelectItem key={item.key} id={item.key}>
              {item.label}
            </AppSelectItem>
          )}
        </AppSelect>
      </div>
      <div className={styles.gameGrid}>
        {games
          ?.sort((a, b) => (a.name < b.name ? -1 : 1))
          .map((game) => (
            <GameCard key={game.titleID} game={game}></GameCard>
          ))}
      </div>
    </div>
  );
}
