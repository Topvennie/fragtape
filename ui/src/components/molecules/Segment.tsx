import { SegmentedControl, SegmentedControlProps } from "@mantine/core";
import classes from "./segment.module.css"

type Props = SegmentedControlProps

export const Segment = (props: Props) => {
  return <SegmentedControl
    withItemsBorders={false}
    color="primary.6"
    styles={{
      root: {
        background: "var(--mantine-color-background-9)",
      },
    }}
    classNames={classes}
    {...props}
  />
}
