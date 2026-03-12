import { isResponseNot200Error } from "@/lib/api/query"
import { useSettingUserFaceitConnect, useSettingUserFaceitDisconnect, useSettingUserGet } from "@/lib/api/setting_user"
import { getErrorMessage } from "@/lib/utils"
import { Button, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { useEffect, useState } from "react"
import { LuCircleCheckBig, LuTrash2, LuTriangleAlert } from "react-icons/lu"
import { Alert } from "../atoms/Alert"
import { ConnectedBadge } from "../atoms/ConnectedBadge"
import { DisconnectedBadge } from "../atoms/DisconnectedBadge"
import { Section } from "../atoms/Page"

export const SettingFaceit = () => {
  const { data: setting } = useSettingUserGet()

  return (
    <Section
      title="Faceit Integration"
      rightSection={setting?.connectedFaceit
        ? <ConnectedBadge />
        : <DisconnectedBadge />
      }
    >
      {setting?.connectedFaceit
        ? <Connected />
        : <Disconnected />
      }
    </Section>
  )
}

const Connected = () => {
  const [submitting, setSubmitting] = useState(false)
  const faceitDisconnect = useSettingUserFaceitDisconnect()

  const handleSubmit = () => {
    setSubmitting(true)

    faceitDisconnect.mutateAsync(undefined, {
      onSuccess: () => notifications.show({ message: "Faceit disconnected" }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => setSubmitting(false),
    })
  }

  return (
    <Stack>
      <Alert
        color="green"
        border={false}
        icon={<LuCircleCheckBig className="size-4 text-green-500" />}
        className="text-xs whitespace-pre-wrap"
      >
        {`Faceit integration is live.\nNew matches will automatically appear in your overview shortly after you finish playing.`}
      </Alert>
      <div className="ml-auto">
        <Button onClick={handleSubmit} color="red" c="black" size="xs" leftSection={<LuTrash2 className="size-4" />} loading={submitting}>
          Disconnect Faceit
        </Button>
      </div>
    </Stack>
  )
}

const Disconnected = () => {
  const [submitting, setSubmitting] = useState(false)
  const [notFound, setNotFound] = useState(false)
  const faceitConnect = useSettingUserFaceitConnect()

  const handleSubmit = () => {
    setSubmitting(true)
    const id = notifications.show({ message: "Connecting...", autoClose: false, loading: true })

    faceitConnect.mutate(undefined, {
      onSuccess: () => notifications.update({ id, message: "Faceit connected", autoClose: true, loading: false }),
      onError: async (error) => {
        if (isResponseNot200Error(error) && error.response.status === 404) {
          setNotFound(true)
          notifications.hide(id)
          return
        }

        const msg = await getErrorMessage(error)
        notifications.update({ id, color: "red", message: msg, autoClose: true, loading: false })
      },
      onSettled: () => setSubmitting(false),
    })
  }

  useEffect(() => {
    if (notFound) setTimeout(() => setNotFound(false), 10000)
  }, [notFound])

  return (
    <Stack>
      {notFound && (
        <Alert
          color="red"
          icon={<LuTriangleAlert className="size-6 text-red-500" />}
          className="text-xs"
        >
          <p>No Faceit acccount linked to this Steam account.</p>
          <p>Double check your Steam account.</p>
        </Alert>
      )}
      <p className="text-white">Connect your Faceit account to enable automatic match tracking for your Faceit matches</p>
      <div className="ml-auto">
        <Button onClick={handleSubmit} loading={submitting} disabled={notFound}>
          Link Faceit account
        </Button>
      </div>
    </Stack>
  )
}
