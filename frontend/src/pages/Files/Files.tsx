import { useFiles } from "../../hooks/useFiles";

export default function Files() {
  const { data, isLoading, error } = useFiles();

  return (
    <div>
      <div>
        {isLoading && "loading"} {`${error}`}
      </div>
      <div>{data && data.map((e) => <div>{JSON.stringify(e)}</div>)}</div>
    </div>
  );
}
