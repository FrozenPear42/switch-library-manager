import { IconCheck, IconCross, IconX } from "@tabler/icons-react";
import { main } from "../../../wailsjs/go/models";
import styles from "./CatalogGameCard.module.css";
import classNames from "classnames";

export function CatalogGameCard({ data }: { data: main.SwitchTitle }) {
  const lastUpdate =
    data.versions.length > 0
      ? {
          releaseDate: data.versions[data.versions.length - 1].releaseDate,
          version: data.versions[data.versions.length - 1].version.toFixed(0),
        }
      : { releaseDate: data.releaseDate, version: data.version };

  return (
    <div className={styles.card}>
      <img src={data.banner} className={styles.backgroundBanner} />
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
            <div
              className={classNames(
                styles.libraryInfo,
                data.inLibrary ? styles.success : styles.fail
              )}
            >
              {data.inLibrary ? (
                <>
                  <IconCheck />
                  In Library
                </>
              ) : (
                <>
                  <IconX />
                  Not in Library
                </>
              )}
            </div>
          </div>
        </div>
        <div>{data.intro}</div>
        {/* <div>{data.description}</div> */}
        <div>
          {data.DLCs.length} DLCs...
          {data.DLCs.map((dlc) => (
            <div key={dlc.titleID}>
              <div>
                DLC:
                {dlc.name} {dlc.version} {dlc.region} {dlc.titleID} {dlc.intro}{" "}
                {dlc.description}
              </div>
              {/* <img src={dlc.banner} height={100} />
              <img src={dlc.icon} height={100} /> */}
            </div>
          ))}
        </div>
        <div>
          {data.versions.length + 1} versions...
          {data.versions.map((version) => (
            <div key={version.version}>
              <div>
                {version.version} {version.releaseDate}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
