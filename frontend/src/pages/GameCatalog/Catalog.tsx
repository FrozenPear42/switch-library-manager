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
          .filter((e) => e.titleID != "" && e.name != "")
          // .filter((e) => e.versions.length > 0)
          .slice(0, 100)
          .map((e) => {
            return <CatalogGameCard data={e} key={e.titleID} />;
          })}
      </div>
    </div>
  );
}
