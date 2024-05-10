import { main } from "../../../wailsjs/go/models";
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

  return (
    <div className={styles.card}>
      <img src={game.icon} height={100} className={styles.image} />
      <div className={styles.details}>
        <div className={styles.title}>{game.name}</div>
        <div className={styles.version}>
          {game.version}{" "}
          <span className={styles.additionalInfo}>(v{game.version})</span>
        </div>
        <div className={styles.ids}>
          {game.titleID}{" "}
          <span className={styles.additionalInfo}>({game.region})</span>
        </div>
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

        {/* 
      <div>29 DLCs (3 missing)</div>
      <div>30 files...</div> */}
      </div>
    </div>
  );
}
