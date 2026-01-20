import { API } from "./api";

export interface Highlight {
  id: number;
  title: string;
  round: number;
  durationS: number;
}

export const convertHighlight = (h: API.Highlight): Highlight => {
  return {
    id: h.id,
    title: h.title,
    round: h.round,
    durationS: h.duration_s,
  }
}
