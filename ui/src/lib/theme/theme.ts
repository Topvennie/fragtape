import { createTheme } from "@mantine/core";

export const theme = createTheme({
  fontFamily: "Inter, sans-serif",
  colors: {
    primary: [
      "#e4fef8",
      "#d4f7ef",
      "#adecdd",
      "#7fe0c9",
      "#60d7bb",
      "#48d2b1",
      "#39cfac",
      "#27b796",
      "#17a384",
      "#008d71"
    ],
    secondary: [
      "#f6edff",
      "#dcc9f3",
      "#caaeeb",
      "#ac81e0",
      "#945cd7",
      "#8444d1",
      "#7c38cf",
      "#6b2ab8",
      "#5f25a5",
      "#521d91"
    ],
    background: [
      "#f7f8f9",
      "#e7e7e7",
      "#cccccd",
      "#aeb1b4",
      "#94999f",
      "#838a92",
      "#7a838d",
      "#67717b",
      "#5a646f",
      "#4a5763"
    ],
  },
  primaryColor: "primary",
  primaryShade: 3,
  cursorType: "pointer",
  breakpoints: {
    xs: "36em",
    sm: "40em",
    md: "48em",
    lg: "64em",
    xl: "80em",
    xxl: "96em",
    xxxl: "142em",
    xxxxl: "172em",
  },
});
