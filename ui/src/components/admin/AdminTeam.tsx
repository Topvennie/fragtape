import { useUserAdmin, useUserCreateAdmin, useUserDeleteAdmin, useUserFiltered } from "@/lib/api/user"
import { useAuth } from "@/lib/hooks/useAuth"
import { User } from "@/lib/types/user"
import { getErrorMessage } from "@/lib/utils"
import { ActionIcon, Button, Center } from "@mantine/core"
import { useDebouncedValue, useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { LuPlus, LuTrash2, LuUserRoundPlus } from "react-icons/lu"
import useInfiniteScroll from "react-infinite-scroll-hook"
import { BottomOfPage } from "../atoms/ButtomOfPage"
import { Card } from "../atoms/Card"
import { ModalCenter } from "../atoms/ModalCenter"
import { Title } from "../atoms/Title"
import { FragtapeIcon } from "../icons/FragtapeIcon"
import { Search } from "../molecules/Search"
import { UserList } from "../user/UserList"

export const AdminTeam = () => {
  const { user } = useAuth()
  const { data: admins, isLoading } = useUserAdmin()

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
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Title order={3}>Admin Team</Title>
          <Button onClick={open} leftSection={<LuPlus />} loading={isLoading}>
            Add Admin
          </Button>
        </div>
        <Card>
          {isLoading ? (
            <Center>
              <FragtapeIcon animated className="size-12 text-(--mantine-color-primary-6)" />
            </Center>
          ) : (
            <UserList
              users={admins ?? []}
              leftSection={(u) => {
                if (u.id === user?.id) return <p className="text-secondary">(You)</p>

                return (
                  <ActionIcon onClick={() => handleDelete(u)} variant="subtle" color="muted" loading={deleting}>
                    <LuTrash2 />
                  </ActionIcon>
                )
              }}
            />
          )}
        </Card>
      </div>
      <ModalCenter
        withCloseButton={false}
        opened={opened}
        onClose={close}
      >
        <AddAdmin />
      </ModalCenter>
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
    <div>
      <div className="flex gap-4 items-center py-2 sticky top-0 z-50 bg-(--mantine-color-background-8)">
        <Search
          placeholder="Filter by username..."
          value={name}
          onChange={e => setName(e.target.value)}
          className="grow"
        />
        <p className="text-secondary">{`${result.total} users`}</p>
      </div>
      <UserList
        users={users ?? []}
        leftSection={(u) => (
          <ActionIcon variant="subtle" color="muted" onClick={() => handleCreate(u)} loading={creating.includes(u.id)}>
            <LuUserRoundPlus />
          </ActionIcon>
        )}
        isLoading={isLoading}
      />
      <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
    </div>
  )
}
