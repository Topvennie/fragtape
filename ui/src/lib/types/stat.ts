import { API } from "./api";

export enum Result {
  Win = "win",
  Loss = "loss",
  Tie = "tie",
}
export const resultString: Record<Result, string> = {
  [Result.Win]: "Win",
  [Result.Loss]: "Loss",
  [Result.Tie]: "Tie",
}

export enum Team {
  CT = "ct",
  T = "t",
}

export interface Stat {
  result: Result;
  startTeam: Team;
  kills: number;
  assists: number;
  deaths: number;
}

export const convertStat = (s: API.Stat): Stat => {
  return {
    result: s.result as Result,
    startTeam: s.start_team as Team,
    kills: s.kills,
    assists: s.assists,
    deaths: s.deaths,
  }
}
