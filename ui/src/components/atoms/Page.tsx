import { cn } from "@/lib/utils";
import { Group, Stack, StackProps } from "@mantine/core";
import { Title } from "./Title";
import { ReactNode } from "react";
import { Card } from "./Card";

type PageProps = StackProps

export const Page = ({ className, ...props }: PageProps) => {
  return <Stack gap="xl" p={{ base: "sm", md: "lg" }} className={cn("flex-1", className)} {...props} />
}

type PageTitleProps = {
  title: string;
  description?: string;
  rightSection?: ReactNode;
} & StackProps

export const PageTitle = ({ title, description, rightSection, ...props }: PageTitleProps) => {
  return (
    <Group justify="space-between">
      <Stack gap="xs" {...props}>
        <Title order={1}>{title}</Title>
        <p className="text-secondary">{description}</p>
      </Stack>
      {rightSection}
    </Group>
  )
}

type SectionProps = {
  title?: string;
  rightSection?: ReactNode;
  card?: boolean;
} & StackProps

export const Section = ({ title, rightSection, card = true, className, children, ...props }: SectionProps) => {
  return (
    <Stack className={cn("", className)} {...props}>
      <Group justify="space-between">
        <Title order={2}>{title}</Title>
        {rightSection}
      </Group>

      {card
        ? <Card>{children}</Card>
        : children
      }
    </Stack>
  )
}
