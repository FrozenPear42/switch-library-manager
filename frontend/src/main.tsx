import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./styles/global.css";
import { QueryClient, QueryClientProvider } from "react-query";

const container = document.getElementById("root");

const root = createRoot(container!);

const queryClient = new QueryClient();

root.render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <main>
        <App />
      </main>
    </QueryClientProvider>
  </React.StrictMode>
);
