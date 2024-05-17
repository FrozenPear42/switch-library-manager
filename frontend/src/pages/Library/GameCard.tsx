import {
  IconAlertCircleFilled,
  IconAlertTriangle,
  IconDatabase,
  IconExclamationCircle,
  IconInfoCircle,
  IconInfoCircleFilled,
  IconNetwork,
  IconWorld,
} from "@tabler/icons-react";
import { main } from "../../../wailsjs/go/models";
import { Card } from "../../components/Card/Card";
import styles from "./GameCard.module.css";

export type LocalGameInfo = {
  image: string;
  title: string;
  titleID: string;
  region: string;
  fileType: string;
  updateNumber: number;
  version: string;
};

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

  return (
    <Card coverImage={game.banner} className={styles.card}>
      <div className={styles.titleBox}>
        <img src={game.icon} className={styles.icon} />
        <div>
          <div className={styles.title}>{game.name}</div>
          <div>
            <IconInfoCircleFilled />
          </div>
          <div className={styles.ids}>
            {game.titleID}{" "}
            <span className={styles.additionalInfo}>({game.region})</span>
          </div>

          <div>Base game version: -</div>
          <div>
            <IconAlertCircleFilled />
            Update version: <IconDatabase />
            {`1.0.1`} <IconWorld /> {`1.0.2`}
          </div>
          <div>
            DLCs:{" "}
            {allDLCs.length > 0
              ? `${allDLCs.length - missingDLCs.length}/${allDLCs.length}`
              : "No DLCs available"}
          </div>
          <div>Files: 29 NSP, 1XCI</div>
        </div>
      </div>
      <div></div>

      {/* <div>
        Issues
        <div>
          {missingDLCs.map((dlc) => (
            <div>{`${dlc.titleID} ${dlc.name}`}</div>
          ))}
        </div>
      </div>
      <div>
        Files
        <div>{game.files[0]?.filePath}</div>
      </div> */}
      {/* 
      <div className={styles.details}>
        <div>
          <div>Base game</div>
          {game.files.map((file) => (
            <div>
              {file.readableVersion} {file.filePath}
            </div>
          ))}
        </div>
        <div>
          <div>Updates</div>
          {Object.entries(game.updates).map(([id, update]) => (
            <div>
              {id}
              {update.files.map((file) => (
                <div key={file.fileID}>
                  {file.readableVersion} {file.filePath}
                </div>
              ))}
            </div>
          ))}
        </div>
        <div>
          <div>DLCs</div>
          {Object.entries(game.dlcs).map(([id, dlc]) => (
            <div>
              {id} {dlc.name}
              {dlc.files?.map((file) => (
                <div>
                  {file.fileVersion} {file.filePath}
                </div>
              ))}
            </div>
          ))}
        </div>
      </div> */}
    </Card>
  );
}
