import { LinkButton } from "@/components/atoms/LinkButton";
import { UserIcon } from "@/components/atoms/UserIcon";
import { FragtapeIcon } from "@/components/icons/FragtapeIcon";
import { useAuth } from "@/lib/hooks/useAuth";
import { ActionIcon, AppShell, Burger, Container, Group, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LinkProps, useNavigate } from "@tanstack/react-router";
import { ComponentProps } from "react";
import { LuLogOut } from "react-icons/lu";

type Props = ComponentProps<"div">

type Route = {
  title: string;
  link: LinkProps;
};

const routes: Route[] = [
  {
    title: "Overview",
    link: { to: "/" },
  },
];

const NavLink = ({ route, close }: { route: Route, close?: () => void }) => {
  return (
    <div onClick={close}>
      <LinkButton
        to={route.link.to}
        activeProps={{ variant: "subtle", c: "white" }}
        variant="subtle"
        size="md"
        radius="md"
        c="muted"
      >
        {route.title}
      </LinkButton>
    </div>
  );
};

export const NavLayout = ({ children }: Props) => {
  const { user, logout } = useAuth()

  const [opened, { close, toggle }] = useDisclosure();
  const navigate = useNavigate()

  const handleHome = () => {
    navigate({ to: "/" })
  }

  return (
    <AppShell
      header={{ height: 60 }}
      footer={{ height: 60 }}
      navbar={{ width: 300, breakpoint: "lg", collapsed: { desktop: true, mobile: !opened } }}
      className="max-h-screen overflow-hidden"
    >
      <AppShell.Header px="md" withBorder={false} bg="background.9">
        <Group h="100%" justify="space-between">
          <Group gap="xs" visibleFrom="lg">
            <ActionIcon onClick={handleHome} size="xl" variant="subtle">
              <FragtapeIcon className="size-8 text-(--mantine-color-primary-6)" />
            </ActionIcon>
            <p className="font-bold text-xl text-primary mr-4">Fragtape</p>
            {routes.map(r => <NavLink key={r.title} route={r} />)}
          </Group>
          <Burger
            color="white"
            opened={opened}
            onClick={toggle}
            hiddenFrom="lg"
          />
          <Group>
            <p className="text-primary font-semibold">{user!.displayName}</p>
            <UserIcon user={user!} />
            <ActionIcon onClick={logout} color="red" size="xl" variant="subtle">
              <LuLogOut />
            </ActionIcon>
          </Group>
        </Group>
      </AppShell.Header>
      <AppShell.Navbar p="md" bg="background.9">
        <Stack>
          {routes.map(r => <NavLink key={r.title} route={r} close={close} />)}
        </Stack>
      </AppShell.Navbar>
      <AppShell.Main>
        <Container fluid className="container mx-auto pt-20">
          {children}
        </Container>
      </AppShell.Main>
    </AppShell>
  )
}
