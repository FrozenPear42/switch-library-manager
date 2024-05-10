import { useFiles } from "../../hooks/useFiles";
import {useLibrary} from "../../hooks/useLibrary";

export default function Files() {
  const { data, isLoading, error } = useLibrary();

  return (
    <div>
      <div>
        {isLoading && "loading"} {`${error}`}
      </div>
      <div>{data && data.map((e) => <div>{JSON.stringify(e)}</div>)}</div>
    </div>
  );
}
