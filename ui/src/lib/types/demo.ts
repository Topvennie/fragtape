import { API } from "./api";
import { convertHighlight, Highlight } from "./highlight";
import { convertStat, Stat } from "./stat";
import { convertStatsDemo, StatsDemo } from "./stats_demo";
import { convertUser, User } from "./user";

export enum DemoSource {
  Manual = "manual",
  Steam = "steam",
  Faceit = "faceit",
}
export const demoSourceString: Record<DemoSource, string> = {
  [DemoSource.Manual]: "Manual",
  [DemoSource.Steam]: "Matchmaking",
  [DemoSource.Faceit]: "Faceit",
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
  players: DemoPlayer[];
  stats: StatsDemo;
  createdAt: Date;
  statusUpdatedAt: Date;
}

export interface DemoPlayer {
  user: User;
  stat: Stat;
  highlights: Highlight[];
}

export const convertDemo = (d: API.Demo): Demo => {
  return {
    id: d.id,
    source: d.source as DemoSource,
    status: d.status as DemoStatus,
    players: d.players.map(convertDemoPlayer),
    stats: convertStatsDemo(d.stats),
    createdAt: new Date(d.created_at),
    statusUpdatedAt: new Date(d.status_updated_at),
  }
}

export const convertDemos = (d: API.Demo[]): Demo[] => {
  return d.map(convertDemo)
}

export const convertDemoPlayer = (p: API.DemoPlayer): DemoPlayer => {
  return {
    user: convertUser(p.user),
    stat: convertStat(p.stat),
    highlights: (p.highlights ?? []).map(convertHighlight)
  }
}
