import { useCatalog } from "../../hooks/useCatalog";
import { CatalogGameCard } from "./CatalogGameCard";

import styles from "./Catalog.module.css";

export default function Catalog() {
  const { data, isLoading, error } = useCatalog();

  return (
    <div>
      <div>Filters</div>
      <div>
        {isLoading && "loading"} {error}
      </div>
      <div className={styles.gameList}>
        {data
          .slice(0, 100)
          .filter((e) => e.titleID != "" && e.name != "")
          .map((e) => {
            return <CatalogGameCard data={e} key={e.titleID} />;
          })}
      </div>
    </div>
  );
}
