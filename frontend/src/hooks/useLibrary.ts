import { useQuery } from "react-query";
import {LoadLibraryGames} from "../../wailsjs/go/main/App";

export const useLibrary = () => {
  const { data, isLoading, error } = useQuery(
    "library",
    async () => await LoadLibraryGames()
  );

  return {
    data,
    isLoading,
    error,
  };
};
