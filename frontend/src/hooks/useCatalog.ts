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

export const useCatalog = (
  page: number,
  pageSize: number,
  filters?: CatalogFilters
): HookReturnType => {
  const { data, isLoading, error } = useQuery(
    ["catalog", filters, page, pageSize],
    async () =>
      await LoadCatalog({
        cursor: page * pageSize,
        limit: pageSize,
        region: [],
        sortBy: "name",
      }),
    { keepPreviousData: true }
  );

  return {
    data,
    isLoading,
    error,
  };
};
