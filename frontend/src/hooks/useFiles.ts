import { useQuery } from "react-query";
import { LoadLibraryFiles } from "../../wailsjs/go/main/App";

export const useFiles = () => {
  const { data, isLoading, error } = useQuery(
    "files",
    async () => await LoadLibraryFiles()
  );

  return {
    data,
    isLoading,
    error,
  };
};
