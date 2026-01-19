import { API } from "./api";

export interface StatsDemo {
  map: string;
  roundsCt: number;
  roundsT: number;
}

export const convertStatsDemo = (s: API.StatsDemo): StatsDemo => {
  return {
    map: s.map,
    roundsCt: s.rounds_ct,
    roundsT: s.rounds_t,
  }
}
