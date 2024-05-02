type LocalGameInfo = {
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
  return <div>0</div>;
}
