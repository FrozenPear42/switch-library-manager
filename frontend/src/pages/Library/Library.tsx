import GameCard, { LocalGameInfo } from "./GameCard";
import styles from "./Library.module.css";

const mockGame: LocalGameInfo = {
  fileType: "NSP",
  image:
    "https://img-eshop.cdn.nintendo.net/i/385efb08704498f73482879bdbe87ea7220433cb9126cd91ae262ee68c0b510a.jpg",
  region: "US",
  title: "Atelier Ryza: Ever Darkness & the Secret Hideout",
  titleID: "0100d1900ec80000",
  updateNumber: 524288,
  version: "1.0.8",
};

export default function Library() {
  return (
    <div>
      Library
      <div>
        <button>Organize Library...</button>
      </div>
      <div>Filters</div>
      <div>
        <GameCard gameInfo={mockGame}></GameCard>
      </div>
    </div>
  );
}
