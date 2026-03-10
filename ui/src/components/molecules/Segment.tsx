import { SegmentedControl, SegmentedControlProps, Stack } from "@mantine/core";
import classes from "./segment.module.css";

type Props = {
  label?: string;
} & SegmentedControlProps

export const Segment = ({ label, ...props }: Props) => {
  return (
    <Stack gap={2}>
      <p className="text-white font-semibold text-sm">{label}</p>
      <SegmentedControl
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

    </Stack>
  )
}
