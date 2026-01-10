import { API } from "./api";

export enum DemoSource {
  Manual = "manual",
  Steam = "steam",
  Faceit = "faceit",
}

export enum DemoStatus {
  QueuedParse = "queued_parse",
  Parsing = "parsing",
  QueuedRender = "queued_render",
  Rendering = "rendering",
  Completed = "completed",
  Failed = "failed",
}

export interface Demo {
  id: number;
  source: DemoSource;
  status: DemoStatus;
  createdAt: Date;
  statusUpdatedAt: Date;
}

export const convertDemo = (d: API.Demo): Demo => {
  return {
    id: d.id,
    source: d.source as DemoSource,
    status: d.status as DemoStatus,
    createdAt: new Date(d.created_at),
    statusUpdatedAt: new Date(d.status_updated_at),
  }
}

export const convertDemos = (d: API.Demo[]): Demo[] => {
  return d.map(convertDemo)
}
