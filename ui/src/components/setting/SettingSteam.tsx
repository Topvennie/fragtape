import { useSettingUserGet, useSettingUserSteamConnect, useSettingUserSteamDisconnect } from "@/lib/api/setting_user"
import { settingUserSteamSchema, SettingUserSteamSchema } from "@/lib/types/setting_user"
import { getErrorMessage } from "@/lib/utils"
import { Button, Stack } from "@mantine/core"
import { useForm } from "@mantine/form"
import { notifications } from "@mantine/notifications"
import { zod4Resolver } from "mantine-form-zod-resolver"
import { useState } from "react"
import { LuArrowRight, LuCircleCheckBig, LuInfo, LuTrash2 } from "react-icons/lu"
import { Alert } from "../atoms/Alert"
import { Section } from "../atoms/Page"
import { Ping } from "../atoms/Ping"
import { TextInput } from "../atoms/TextInput"

export const SettingSteam = () => {
  const { data: setting } = useSettingUserGet()

  return (
    <Section
      title="Steam Integration"
      rightSection={setting?.connectedSteam
        ? <ConnectedBadge />
        : <DisconnectedBadge />
      }
    >
      {setting?.connectedSteam
        ? <Connected />
        : <Disconnected />
      }
    </Section>
  )
}

const Connected = () => {
  const [submitting, setSubmitting] = useState(false)
  const steamDisconnect = useSettingUserSteamDisconnect()

  const handleSubmit = () => {
    setSubmitting(true)

    steamDisconnect.mutateAsync(undefined, {
      onSuccess: () => notifications.show({ message: "Steam disconnected" }),
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
        icon={<LuCircleCheckBig className="size-4 text-green-500" />}
        className="text-xs whitespace-pre-wrap"
      >
        {`Steam integration is live.\nNew matches will automatically appear in your overview shortly after you finish playing.`}
      </Alert>
      <div className="ml-auto">
        <Button onClick={handleSubmit} color="red" c="black" size="xs" leftSection={<LuTrash2 className="size-4" />} loading={submitting}>
          Disconnect Steam
        </Button>
      </div>
    </Stack>
  )
}

const ConnectedBadge = () => {
  return (
    <div className="rounded-full bg-(--mantine-color-green-light) px-3 py-1 flex items-center gap-1 text-green-500 text-xs">
      <Ping healthy />
      Connected
    </div>
  )
}

const Disconnected = () => {
  const [submitting, setSubmitting] = useState(false)
  const steamConnect = useSettingUserSteamConnect()

  const form = useForm<SettingUserSteamSchema>({
    initialValues: {
      match_token: "",
      authentication_token: "",
    },
    validate: zod4Resolver(settingUserSteamSchema),
  })

  const handleSubmit = () => {
    if (form.validate().hasErrors) {
      return
    }

    setSubmitting(true)
    const id = notifications.show({ message: "Testing Connection...", autoClose: false, loading: true })

    steamConnect.mutate(form.getValues(), {
      onSuccess: () => notifications.update({ id, message: "Steam connected", autoClose: true, loading: false }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.update({ id, color: "red", message: msg, autoClose: true, loading: false })
      },
      onSettled: () => setSubmitting(false),
    })
  }

  return (
    <Stack>
      <Alert
        color="blue"
        icon={<LuInfo className="size-6 text-white" />}
        className="text-xs"
      >
        <Stack gap={0}>
          <p>
            <span>To enable automatic match tracking, you need to provide your Steam Match Token and Authentication Token. </span>
            <a
              href="https://help.steampowered.com/en/wizard/HelpWithGameIssue/?appid=730&issueid=128"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-(--mantine-color-primary-6)"
            >
              Find your Steam tokens here
              <LuArrowRight />
            </a>
          </p>
          <p>
            <span>More information can be found in the </span>
            <a
              href="https://leetify.com/blog/share-codes/"
              target="_blank"
              rel="noopener noreferrer"
              className="text-(--mantine-color-primary-6)"
            >
              Leetify blog post.
            </a>
          </p>
        </Stack>
      </Alert>
      <div>
        <p className="text-white">Steam Match Token</p>
        <TextInput placeholder="CSGO-xxxxx-xxxxx-xxxxx-xxxxx-xxxxx" {...form.getInputProps("match_token")} />
      </div>
      <div>
        <p className="text-white">Steam Authentication Token</p>
        <TextInput placeholder="AAAA-BBBBB-CCCC" {...form.getInputProps("authentication_token")} />
      </div>
      <div className="ml-auto">
        <Button onClick={handleSubmit} loading={submitting}>
          Save Steam Settings
        </Button>
      </div>
    </Stack>
  )
}

const DisconnectedBadge = () => {
  return (
    <div className="rounded-full bg-gray-400/10 px-3 py-1 flex items-center gap-1 text-gray-400 text-xs">
      Not connected
    </div>
  )
}
