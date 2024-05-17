import { useQuery } from "react-query";
import { LoadLibraryGames } from "../../wailsjs/go/main/App";

export const useLibrary = () => {
  const { data, isLoading, error } = useQuery("library", async () => {
    const games = await LoadLibraryGames();
    console.log(games);
    return games;
  });

  return {
    data,
    isLoading,
    error,
  };
};
