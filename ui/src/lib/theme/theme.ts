import { colorsTuple, createTheme } from "@mantine/core";

export const theme = createTheme({
  fontFamily: "Inter, sans-serif",
  white: "#E6EDF3",
  colors: {
    primary: [
      "#defeff",
      "#caf8ff",
      "#99edff",
      "#64e3ff",
      "#3ddafe",
      "#25d5fe",
      "#00d1ff",
      "#00bae4",
      "#00a6cc",
      "#0090b4"
    ],
    background: [
      "#f1f4f9",
      "#e1e5ec",
      "#bec8da",
      "#98aac8",
      "#7990ba",
      "#6580b1",
      "#5a78ae",
      "#4a6699",
      "#121A27",
      "#0f1724"
    ],
    steam: colorsTuple("#171a21"),
  },
  autoContrast: true,
  primaryColor: "primary",
  primaryShade: 6,
  cursorType: "pointer",
  defaultRadius: "md",
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
