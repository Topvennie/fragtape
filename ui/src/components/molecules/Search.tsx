import { TextInput, TextInputProps } from "@mantine/core"
import { LuSearch } from "react-icons/lu"

type Props = TextInputProps

export const Search = (props: Props) => {
  return (
    <TextInput
      leftSection={<LuSearch />}
      radius="md"
      c="white"
      autoComplete="off"
      styles={{
        input: {
          background: "var(--mantine-color-background-8)",
          borderColor: "var(--mantine-color-background-7)",
          color: "white",
        }
      }}
      {...props}
    />
  )
}
