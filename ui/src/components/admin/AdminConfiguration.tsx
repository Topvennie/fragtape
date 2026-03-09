import { convertSettingGlobalSchema, SettingGlobal, settingGlobalSchema, SettingGlobalSchema } from "@/lib/types/setting_global"
import { ActionIcon, Button, Center, Group, Stack, Switch } from "@mantine/core"
import { Block } from "@tanstack/react-router"
import { LuSave } from "react-icons/lu"

import { useSettingGlobalGet, useSettingGlobalUpdate } from "@/lib/api/setting_global"
import { getErrorMessage } from "@/lib/utils"
import { useForm } from "@mantine/form"
import { notifications } from "@mantine/notifications"
import { zod4Resolver } from "mantine-form-zod-resolver"
import { ReactNode, useEffect, useState } from "react"
import { ModalCenter } from "../atoms/Modal"
import { Section } from "../atoms/Page"
import { TextInput } from "../atoms/TextInput"

type Keys = keyof SettingGlobal
type Values = SettingGlobal[Keys]

export const AdminConfiguration = () => {
  const { data: originalSetting } = useSettingGlobalGet()
  const settingUpdate = useSettingGlobalUpdate()

  const [updating, setUpdating] = useState(false)

  const form = useForm<SettingGlobalSchema>({
    initialValues: convertSettingGlobalSchema(originalSetting!),
    validate: zod4Resolver(settingGlobalSchema)
  })
  useEffect(() => {
    if (!originalSetting) return
    const v = convertSettingGlobalSchema(originalSetting)

    form.setInitialValues(v)
    form.setValues(v)
    form.resetDirty()
  }, [originalSetting]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleChange = (k: Keys, value: Values) => {
    form.setFieldValue(k, value)
  }

  const handleUpdate = () => {
    if (form.validate().hasErrors) {
      return
    }

    setUpdating(true)

    settingUpdate.mutate(form.getValues(), {
      onSuccess: () => {
        notifications.show({ message: "Settings saved" })

        const v = convertSettingGlobalSchema(form.getValues())

        form.setInitialValues(v)
        form.setValues(v)
        form.resetDirty()
      },
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      },
      onSettled: () => setUpdating(false),
    })
  }

  return (
    <>
      <Section
        title="Configuration"
        rightSection={(
          <ActionIcon onClick={handleUpdate} variant="subtle" loading={updating} disabled={!form.isDirty()}>
            <LuSave />
          </ActionIcon>
        )}
      >
        <Stack gap="xl">
          <SettingBoolean
            title="Allow Demo Uploads"
            description="Enable users to manually upload demo files"
            checked={form.getValues().demoUpload}
            onChange={checked => handleChange("demoUpload", checked)}
          />
          <SettingBoolean
            title="Custom Highlight Criteria"
            description="Allow regular users to modify their personal highlight generation rules"
            checked={form.getValues().customCriteria}
            onChange={checked => handleChange("customCriteria", checked)}
          />
          <SettingBoolean
            title="Enable Chat Commands"
            description="Listen for specific phrases in game chat to mark rounds for highlights"
            checked={form.getValues().chatCommand}
            onChange={checked => handleChange("chatCommand", checked)}
          />
          <SettingString
            title="Trigger Phrase"
            description="The text command users must type (e.g. fragtape)"
            value={form.getValues().chatTrigger}
            onChange={value => handleChange("chatTrigger", value)}
          />
        </Stack>

        <Block
          shouldBlockFn={() => form.isDirty()}
          enableBeforeUnload={() => form.isDirty()}
          withResolver
        >
          {({ status, proceed, reset }) => (
            <ModalCenter opened={status === "blocked"} onClose={() => reset?.()} size="md">
              <SettingConfirm proceed={proceed} reset={reset} />
            </ModalCenter>
          )}
        </Block>
      </Section>
    </>
  )
}

type SettingBooleanProps = {
  title: string;
  description: string;
  checked: boolean;
  onChange: (value: boolean) => void;
}

const SettingBoolean = ({ title, description, checked, onChange }: SettingBooleanProps) => {
  return <Setting
    title={title}
    description={description}
    rightSection={(
      <Switch
        checked={checked}
        onChange={e => onChange(e.target.checked)}
        size="md"
      />
    )}
  />
}

type SettingStringProps = {
  title: string;
  description: string;
  value: string;
  onChange: (value: string) => void;
}

const SettingString = ({ title, description, value, onChange }: SettingStringProps) => {
  return <Setting
    title={title}
    description={description}
    rightSection={(
      <TextInput
        value={value}
        onChange={e => onChange(e.currentTarget.value)}
      />
    )}
  />
}

type SettingProps = {
  title: string;
  description: string;
  rightSection: ReactNode;
}

const Setting = ({ title, description, rightSection }: SettingProps) => {
  return (
    <Group justify="space-between">
      <Stack gap={0}>
        <p className="text-white">{title}</p>
        <p className="text-sm text-secondary">{description}</p>
      </Stack>
      {rightSection}
    </Group>
  )
}

type SettingConfirm = {
  proceed?: () => void;
  reset?: () => void;
}

const SettingConfirm = ({ proceed, reset }: SettingConfirm) => {
  return (
    <Stack>
      <Center className="flex flex-col">
        <p>You have unchanged changes</p>
        <p>Are you sure you want to leave?</p>
      </Center>
      <Group justify="end">
        <Button onClick={reset} variant="outline">
          Go Back
        </Button>
        <Button onClick={proceed}>
          Proceed
        </Button>
      </Group>
    </Stack>
  )
}
