import { IconCheck, IconCross, IconX } from "@tabler/icons-react";
import { main } from "../../../wailsjs/go/models";
import styles from "./CatalogGameCard.module.css";
import classNames from "classnames";
import Collapsible from "../../components/Collapsible/Collapsible";
import { Card } from "../../components/Card/Card";

function LibraryIndicator({ inLibrary }: { inLibrary: boolean }) {
  return (
    <div
      className={classNames(
        styles.libraryInfo,
        inLibrary ? styles.success : styles.fail
      )}
    >
      {inLibrary ? (
        <>
          <IconCheck />
          <span>In the Library</span>
        </>
      ) : (
        <>
          <IconX />
          <span>Not in the Library</span>
        </>
      )}
    </div>
  );
}

function DLCCard({ data }: { data: main.CatalogDLCData }) {
  return (
    <div className={styles.dlcCard}>
      <img src={data.banner} className={styles.dlcBanner} />
      <div className={styles.dlcDetails}>
        <div>{data.name}</div>
        <div>
          {data.titleID}{" "}
          <span className={styles.additionalInfo}>({data.region})</span>
        </div>
        <LibraryIndicator inLibrary={false}></LibraryIndicator>
      </div>
    </div>
  );
}

export function CatalogGameCard({ data }: { data: main.CatalogSwitchGame }) {
  const lastUpdate =
    data.versions.length > 0
      ? {
          releaseDate: data.versions[data.versions.length - 1].releaseDate,
          version: data.versions[data.versions.length - 1].version.toFixed(0),
        }
      : { releaseDate: data.releaseDate, version: data.version };

  return (
    <Card coverImage={data.banner} className={styles.card}>
      <div className={styles.content}>
        <div className={styles.header}>
          <img src={data.icon} className={styles.gameIcon} />
          <div className={styles.details}>
            <div className={styles.title}>{data.name}</div>
            <div>
              {data.titleID}{" "}
              <span className={styles.additionalInfo}>({data.region})</span>
            </div>
            <div>
              Update: {`v${lastUpdate.version}`}{" "}
              <span className={styles.additionalInfo}>
                ({lastUpdate.releaseDate})
              </span>
            </div>
            <LibraryIndicator inLibrary={false}></LibraryIndicator>
          </div>
        </div>
        <div>{data.intro}</div>
        {/* <div>{data.description}</div> */}
        {data.dlcs.length > 0 && (
          <>
            <div className={styles.spacer}></div>
            <div>{data.dlcs.length} DLCs available </div>
            {data.dlcs.length > 1 ? (
              <Collapsible
                startOpened={false}
                openText="Show all..."
                closeText="Hide"
              >
                <div className={styles.dlcList}>
                  {data.dlcs.map((dlc) => (
                    <DLCCard data={dlc} key={dlc.titleID} />
                  ))}
                </div>
              </Collapsible>
            ) : (
              <div className={styles.dlcList}>
                <DLCCard data={data.dlcs[0]} />
              </div>
            )}
          </>
        )}
      </div>
    </Card>
  );
}
