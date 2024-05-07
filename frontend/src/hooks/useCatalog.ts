import { useEffect, useState } from "react";
import { LoadCatalog } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";
import { useQuery } from "react-query";

export type CatalogFilters = {
  name: string | null;
};

type HookReturnType = {
  data: main.CatalogPage | undefined;
  isLoading: boolean;
  error: unknown;
};

export const useCatalog = (filters?: CatalogFilters): HookReturnType => {
  const { data, isLoading, error } = useQuery(
    ["catalog", filters],
    async () =>
      await LoadCatalog({
        cursor: 100,
        limit: 100,
        region: [],
        sortBy: "name",
      })
  );

  return {
    data,
    isLoading,
    error,
  };
};
