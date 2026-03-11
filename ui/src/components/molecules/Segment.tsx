import { SegmentedControl, SegmentedControlProps, Stack, StackProps } from "@mantine/core";
import classes from "./segment.module.css";
import { ComponentProps } from "react";
import { cn } from "@/lib/utils";

type Props = {
  label?: string;
  labelProps?: ComponentProps<"p">;
  segmentProps?: SegmentedControlProps;
} & StackProps

export const Segment = ({
  label,
  labelProps: { className: labelClassName, ...labelProps } = {},
  segmentProps: { className: segmentClassName, ...segmentProps } = { data: [] },
  className,
  ...props
}: Props) => {
  return (
    <Stack gap={2} className={cn("w-fit", className)} {...props}>
      <p className={cn("text-white font-semibold text-sm", labelClassName)} {...labelProps}>{label}</p>
      <SegmentedControl
        withItemsBorders={false}
        color="primary.6"
        styles={{
          root: {
            background: "var(--mantine-color-background-8)",
          },
        }}
        classNames={classes}
        className={segmentClassName}
        {...segmentProps}
      />

    </Stack>
  )
}
