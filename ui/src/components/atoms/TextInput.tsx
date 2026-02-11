import { TextInput as MantineTextInput, TextInputProps } from "@mantine/core"

export const TextInput = (props: TextInputProps) => {
  return (
    <MantineTextInput
      styles={{
        input: {
          background: "var(--mantine-color-background-8)",
          border: "none",
          color: "var(--mantine-color-muted-1)",
        },
      }}
      {...props}
    />
  )
}
