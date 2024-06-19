import {
  IconAlertCircleFilled,
  IconCheck,
  IconDatabase,
  IconExclamationCircle,
  IconWorld,
} from "@tabler/icons-react";
import { main } from "../../../wailsjs/go/models";
import { Card } from "../../components/Card/Card";
import styles from "./GameCard.module.css";
import classNames from "classnames";

type GameCardProps = {
  game: main.LibrarySwitchGame;
};

export default function GameCard({ game }: GameCardProps) {
  // all files list?
  // list of missing dlcs
  // missing update info
  // duplicate file info
  // name, icon, banner, id, version
  // dlc list
  // info if its in library
  // tilte in red/yellow if issues?

  const allDLCs = Object.values(game.dlcs).sort(
    (a, b) => Number.parseInt(a.titleID, 16) - Number.parseInt(b.titleID, 16)
  );
  const missingDLCs = allDLCs.filter((g) => !g.inLibrary);

  const isBaseInLibrary = game.inLibrary;

  const recentCatalogVersion = game.allVersions
    .sort((a, b) => a.version - b.version)
    .at(-1);

  const recentLibraryVersion = Object.values(game.updates)
    .flatMap((u) => u.files)
    .sort((a, b) => a.fileVersion - b.fileVersion)
    .at(-1);

  const files = [
    ...game.files.map((f) => f.filePath),
    ...Object.values(game.dlcs)
      .flatMap((d) => d.files)
      .filter((f) => !!f)
      .map((f) => f.filePath),
    ...Object.values(game.updates)
      .flatMap((u) => u.files)
      .map((f) => f.filePath),
  ];

  const fileTypesMap = files
    .map((fileName) => fileName.split(".").at(-1))
    .map((n) => (n ? n.toUpperCase() : "noext"))
    .reduce(
      (p, n) => ({ ...p, [n]: (p[n] ?? 0) + 1 }),
      {} as Record<string, number>
    );

  const fileTypes = Object.entries(fileTypesMap)
    .map(([k, v]) => ({ type: k, count: v }))
    .sort((a, b) => a.count - b.count);

  return (
    <Card coverImage={game.banner} className={styles.card}>
      <div className={styles.titleBox}>
        <img src={game.icon} className={styles.icon} />
        <div>
          <div className={styles.title}>{game.name}</div>
          {/* {fileTypes.map((f) => (
            <span>{f}</span>
          ))} */}
          <div className={styles.ids}>
            {game.titleID}{" "}
            <span className={styles.additionalInfo}>({game.region})</span>
          </div>

          <div className={styles.details}>
            <div>Base:</div>
            <span
              className={classNames(
                isBaseInLibrary ? styles.colorSuccess : styles.colorError,
                styles.detailsLine
              )}
            >
              {isBaseInLibrary ? (
                <>
                  <IconCheck size={"1rem"} />
                  <span>In library</span>
                </>
              ) : (
                <>
                  <IconExclamationCircle size={"1rem"} />
                  <span>Missing</span>
                </>
              )}
            </span>

            <div>Update:</div>
            <span
              className={classNames(
                game.isRecentUpdateInLibrary
                  ? styles.colorSuccess
                  : styles.colorWarning,
                styles.detailsLine
              )}
            >
              <IconCheck size={"1rem"} />
              <div className={styles.updateInfo}>
                <div className={styles.versionTag}>
                  <IconDatabase size={"1rem"} />
                  {recentLibraryVersion
                    ? `0x${recentLibraryVersion?.fileVersion.toString(16)} ${
                        recentLibraryVersion?.readableVersion
                      }`
                    : "No update in library"}
                </div>
                <div className={styles.versionTag}>
                  <IconWorld size={"1rem"} />
                  {recentCatalogVersion
                    ? `0x${recentCatalogVersion?.version.toString(16)} ${
                        recentCatalogVersion?.releaseDate
                      }`
                    : "No updates available"}
                </div>
              </div>
            </span>

            <div>DLCs:</div>
            <span
              className={classNames(
                missingDLCs.length == 0
                  ? styles.colorSuccess
                  : styles.colorWarning,
                styles.detailsLine
              )}
            >
              {missingDLCs.length == 0 ? (
                <IconCheck size={"1rem"} />
              ) : (
                <IconAlertCircleFilled size={"1rem"} />
              )}

              <span>
                {allDLCs.length > 0
                  ? `${allDLCs.length - missingDLCs.length}/${allDLCs.length}`
                  : "No DLCs available"}
              </span>
            </span>

            <div>Files: </div>
            <span>
              {fileTypes.map((t) => `${t.count} ${t.type}`).join(", ")}
              <button className="button">Show details...</button>
            </span>
          </div>
        </div>
      </div>
    </Card>
  );
}
