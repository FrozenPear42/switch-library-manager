import { useEffect, useState } from "react";
import { LoadCatalog } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";

export type CatalogFilters = {
  name: string | null;
};

type HookReturnType = {
  data: main.SwitchTitle[];
  isLoading: boolean;
  error: string | null;
};

export const useCatalog = (filters?: CatalogFilters): HookReturnType => {
  const [data, setData] = useState<main.SwitchTitle[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetch = async () => {
      setIsLoading(true);
      try {
        const catalog = await LoadCatalog();
        setData(catalog);
      } catch (e) {
        console.error(e);
        setError(`error: ${e}`);
      } finally {
        setIsLoading(false);
      }
    };
    fetch();
  }, [filters]);

  return {
    data,
    isLoading,
    error,
  };
};
