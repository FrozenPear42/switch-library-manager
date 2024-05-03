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
  gameInfo: LocalGameInfo;
};

export default function GameCard({ gameInfo }: GameCardProps) {
  return (
    <div className={styles.card}>
      <img src={gameInfo.image} height={100} className={styles.image} />
      <div className={styles.details}>
        <div className={styles.title}>{gameInfo.title}</div>
        <div className={styles.version}>
          {gameInfo.version}{" "}
          <span className={styles.additionalInfo}>
            (v{gameInfo.updateNumber})
          </span>
        </div>
        <div className={styles.ids}>
          {gameInfo.titleID}{" "}
          <span className={styles.additionalInfo}>({gameInfo.region})</span>
        </div>
        <div>{gameInfo.fileType}</div>
        {/* 
      <div>29 DLCs (3 missing)</div>
      <div>30 files...</div> */}
      </div>
    </div>
  );
}
