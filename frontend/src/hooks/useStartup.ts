import { useEffect, useState, useSyncExternalStore } from "react";
import App from "../App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { RequestStartupProgress } from "../../wailsjs/go/main/App";

import { EventMessage, EventType } from "../model/events";
export type StartupState = {
  running: boolean;
  completed: boolean;
  stageMessage: string;
  stageCurrent: number;
  stageTotal: number;
};

type HookReturnType = {
  //   start: () => void;
  state: StartupState | null;
};

export const useStartup = (): HookReturnType => {
  const [startupState, setStartupState] = useState<StartupState | null>(null);

  useEffect(() => {
    const unsubscribe = EventsOn(
      EventType.StartupProgress,
      (payload: EventMessage) => {
        console.log("received startup event", payload);
        if (payload.type !== EventType.StartupProgress) {
          return;
        }
        setStartupState(<StartupState>{
          completed: payload.data.completed,
          running: payload.data.running,
          stageMessage: payload.data.message,
          stageCurrent: payload.data.current,
          stageTotal: payload.data.total,
        });
      }
    );
    RequestStartupProgress();
    return unsubscribe;
  }, [setStartupState]);

  return {
    state: startupState,
  };
};
