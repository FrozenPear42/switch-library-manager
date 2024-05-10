import { useLibrary } from "../../hooks/useLibrary";
import GameCard from "./GameCard";
import styles from "./Library.module.css";

export default function Library() {
  const { data: games, isLoading, error } = useLibrary();

  return (
    <div>
      <div className={styles.actions}>
        <div>
          <button>Organize Library...</button>
        </div>
        <div>Filters</div>
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
