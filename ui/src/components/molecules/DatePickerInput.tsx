import { cn } from "@/lib/utils";
import { Stack, StackProps } from "@mantine/core";
import { DatePickerInputProps, DatePickerInput as MDatePickerInput, DatePickerType } from "@mantine/dates";
import { ComponentProps } from "react";
import { LuCalendar } from "react-icons/lu";

type Props<T extends DatePickerType> = {
  label?: string;
  labelProps?: ComponentProps<"p">;
  dateProps?: DatePickerInputProps<T>;
} & StackProps;

export const DatePickerInput = <T extends DatePickerType>({
  label,
  labelProps: { className: labelClassName, ...labelProps } = {},
  dateProps: { className: dateClassName, ...dateProps } = {},
  className,
  ...props
}: Props<T>) => {
  return (
    <Stack gap={2} className={cn("w-fit", className)} {...props}>
      <p className={cn("text-white font-semibold text-sm", labelClassName)} {...labelProps}>{label}</p>
      <MDatePickerInput
        leftSection={<LuCalendar />}
        styles={{
          input: {
            background: "var(--mantine-color-background-8)",
            borderColor: "transparent",
            color: "white",
          },
        }}
        className={cn("grow", dateClassName)}
        {...dateProps}
      />
    </Stack>
  )
}
