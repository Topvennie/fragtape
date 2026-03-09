import { LinkButton } from "@/components/atoms/LinkButton";
import { FragtapeIcon } from "@/components/icons/FragtapeIcon";
import { UserIcon } from "@/components/user/UserIcon";
import { useAuth } from "@/lib/hooks/useAuth";
import { ActionIcon, AppShell, Burger, Button, Group, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LinkProps, useNavigate } from "@tanstack/react-router";
import { ComponentProps } from "react";
import { LuLogOut } from "react-icons/lu";

type Props = ComponentProps<"div">

type Route = {
  title: string;
  link: LinkProps;
  admin?: boolean;
};

const allRoutes: Route[] = [
  {
    title: "Overview",
    link: { to: "/" },
  },
  {
    title: "Settings",
    link: { to: "/setting" },
  },
  {
    title: "Admin",
    link: { to: "/admin" },
    admin: true,
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

  const routes = allRoutes.filter(r => !r.admin || user?.admin)

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 300, breakpoint: "md", collapsed: { desktop: true, mobile: !opened } }}
    >
      <AppShell.Header px="md" withBorder={false} bg="background.9">
        <Group h="100%" justify="space-between">
          <Group gap="xs" visibleFrom="md">
            <Button onClick={handleHome} size="md" c="white" variant="subtle" leftSection={<FragtapeIcon className="size-8 text-(--mantine-color-primary-6)" />}>
              Fragtape
            </Button>
            {routes.map(r => <NavLink key={r.title} route={r} />)}
          </Group>
          <Burger
            color="white"
            opened={opened}
            onClick={toggle}
            hiddenFrom="md"
          />
          <Group>
            <p className="text-primary font-semibold hidden md:block">{user!.displayName}</p>
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
      <AppShell.Main className="flex flex-col justify-between">
        <div className="container mx-auto">
          {children}
        </div>
        {/* TODO: Add footer here */}
      </AppShell.Main>

    </AppShell>
  )
}
