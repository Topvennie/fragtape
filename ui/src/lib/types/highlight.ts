import { API } from "./api";

export interface Highlight {
  id: number;
  title: string;
}

export const convertHighlight = (h: API.Highlight): Highlight => {
  return {
    id: h.id,
    title: h.title,
  }
}
