export enum EventType {
  StartupProgress = "startupProgress",
}

export type StartupProgressPayload = {
  completed: boolean;
  running: boolean;
  message: string;
  current: number;
  total: number;
};

export type EventMessage = {
  type: EventType.StartupProgress;
  data: StartupProgressPayload;
};
