import { IconCheck } from "@tabler/icons-react";
import styles from "./NUTPage.module.css";
import Progress from "../../components/Progress/Progress";
import { Card } from "../../components/Card/Card";

type DownloadCardProps = {
  filePath: string;
};

function DownloadCard() {
  return (
    <div>
      <Progress value={50} />
    </div>
  );
}

export default function NUTPage() {
  return (
    <div className={styles.page}>
      <Card>
        <div>
          Server status: <IconCheck /> Running
        </div>
        <div>
          <div>IP: 10.245.13.85</div>
          <div>Port: 9000</div>
        </div>
      </Card>
      <Card></Card>
      <div>
        Activity
        <DownloadCard></DownloadCard>
      </div>
      <div>Log</div>
    </div>
  );
}
