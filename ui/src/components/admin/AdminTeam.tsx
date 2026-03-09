import { useUserAdmin, useUserCreateAdmin, useUserDeleteAdmin, useUserFiltered } from "@/lib/api/user"
import { useAuth } from "@/lib/hooks/useAuth"
import { User } from "@/lib/types/user"
import { getErrorMessage } from "@/lib/utils"
import { ActionIcon, Button, Group, Stack } from "@mantine/core"
import { useDebouncedValue, useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { LuPlus, LuTrash2, LuUserRoundPlus } from "react-icons/lu"
import useInfiniteScroll from "react-infinite-scroll-hook"
import { BottomOfPage } from "../atoms/ButtomOfPage"
import { Modal } from "../atoms/Modal"
import { Section } from "../atoms/Page"
import { Search } from "../molecules/Search"
import { UserList } from "../user/UserList"
import { Loading } from "../atoms/Loading"

export const AdminTeam = () => {
  const { user } = useAuth()
  const { data: admins } = useUserAdmin()

  const [deleting, setDeleting] = useState(false)

  const deleteAdmin = useUserDeleteAdmin()
  const handleDelete = (admin: User) => {
    setDeleting(true)

    deleteAdmin.mutateAsync(admin, {
      onSuccess: () => notifications.show({ message: "Admin removed" }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => setDeleting(false)
    })
  }

  const [opened, { open, close }] = useDisclosure()

  return (
    <>
      <Section
        title="Team"
        rightSection={(
          <Button onClick={open} leftSection={<LuPlus />}>
            Add Admin
          </Button>
        )}
      >
        <UserList
          users={admins ?? []}
          rightSection={(u) => {
            if (u.id === user?.id) return <p className="text-secondary">(You)</p>

            return (
              <ActionIcon onClick={() => handleDelete(u)} variant="subtle" color="muted" loading={deleting}>
                <LuTrash2 />
              </ActionIcon>
            )
          }}
        />
      </Section>

      <Modal
        withCloseButton={false}
        opened={opened}
        onClose={close}
      >
        <AddAdmin />
      </Modal>
    </>
  )
}

const AddAdmin = () => {
  const [name, setName] = useState("")
  const [debouncedName] = useDebouncedValue(name, 200);

  const { result, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } = useUserFiltered({ name: debouncedName, admin: false, real: true });
  const users = result.users.filter(u => !u.admin)

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  const [creating, setCreating] = useState<number[]>([])

  const createAdmin = useUserCreateAdmin()
  const handleCreate = (admin: User) => {
    setCreating(prev => {
      const newUsers = [...prev]
      newUsers.push(admin.id)
      return newUsers
    })

    createAdmin.mutateAsync(admin, {
      onSuccess: () => notifications.show({ message: "Admin added" }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => setCreating(prev => prev.filter(p => p !== admin.id))
    })
  }

  return (
    <Stack>
      <Group className="sticky top-0 z-50">
        <Search
          placeholder="Filter by username..."
          value={name}
          onChange={e => setName(e.target.value)}
          className="grow"
        />
        <p className="text-secondary">{`${result.total} users`}</p>
      </Group>

      <UserList
        users={users ?? []}
        rightSection={(u) => (
          <ActionIcon variant="subtle" color="muted" onClick={() => handleCreate(u)} loading={creating.includes(u.id)}>
            <LuUserRoundPlus />
          </ActionIcon>
        )}
        isLoading={isLoading}
      />

      <Loading isFetchingNextPage={isFetchingNextPage} hasNextPage={hasNextPage} />

      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </Stack>
  )
}
