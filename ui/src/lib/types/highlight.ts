import { API } from "./api";

export interface Highlight {
  id: number;
  title: string;
  round: number;
  kills: number;
  durationS: number;
  generated: boolean;
}

export const convertHighlight = (h: API.Highlight): Highlight => {
  return {
    id: h.id,
    title: h.title,
    round: h.round,
    kills: h.kills,
    durationS: h.duration_s,
    generated: h.generated,
  }
}
