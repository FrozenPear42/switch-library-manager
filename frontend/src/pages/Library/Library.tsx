import {
  IconExclamationCircle,
  IconMessageCircleExclamation,
} from "@tabler/icons-react";
import Spinner from "../../components/Spinner/Spinner";
import { useLibrary } from "../../hooks/useLibrary";
import GameCard from "./GameCard";
import styles from "./Library.module.css";
import { AppSelect, AppSelectItem } from "../../components/Select/Select";
import { useState } from "react";
import { main } from "../../../wailsjs/go/models";
import OrganizeFilesModal from "./OrganizeFilesModal";
import AppTextField from "../../components/TextField/TextField";

export default function Library() {
  const { data: games, isLoading, error } = useLibrary();
  const [sortMode, setSortMode] = useState<"id" | "name" | "region" | "issues">(
    "name"
  );

  const sorters: Record<
    string,
    (a: main.LibrarySwitchGame, b: main.LibrarySwitchGame) => number
  > = {
    id: (a, b) =>
      Number.parseInt(a.titleID, 16) - Number.parseInt(b.titleID, 16),
    name: (a, b) => (a.name < b.name ? -1 : 1),
    region: (a, b) => 0,
    issues: (a, b) => 0,
  };

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
    <>
      <OrganizeFilesModal
        isOpened={true}
        onOpen={console.log}
      ></OrganizeFilesModal>
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
            // defaultSelectedKey={"name"}
            onSelectionChange={(mode) => setSortMode(mode as any)}
          >
            {(item) => (
              <AppSelectItem key={item.key} id={item.key}>
                {item.label}
              </AppSelectItem>
            )}
          </AppSelect>
          <AppTextField label="ID"></AppTextField>
          <AppTextField label="Name"></AppTextField>
        </div>
        <div className={styles.gameGrid}>
          {games?.sort(sorters[sortMode]).map((game) => (
            <GameCard key={game.titleID} game={game}></GameCard>
          ))}
        </div>
      </div>
    </>
  );
}
