import { ActionIcon, Button, Switch } from "@mantine/core"
import { Block } from "@tanstack/react-router"
import { useState } from "react"
import { LuSave } from "react-icons/lu"
import { Card } from "../atoms/Card"
import { ModalCenter } from "../atoms/ModalCenter"
import { Title } from "../atoms/Title"

type SettingType = {
  uploads: boolean;
  customCriteria: boolean;
  chatCommands: boolean;
  chatTrigger: string;
}

export const AdminConfiguration = () => {
  const originalSettings: SettingType = {
    uploads: true,
    customCriteria: false,
    chatCommands: true,
    chatTrigger: "fragtape"
  }
  const [settings, setSettings] = useState<SettingType>({
    uploads: true,
    customCriteria: false,
    chatCommands: true,
    chatTrigger: "fragtape"
  })
  const changed = JSON.stringify(originalSettings) !== JSON.stringify(settings)

  const handleChange = (k: keyof SettingType, checked: boolean) => {
    setSettings(prev => ({ ...prev, [k]: checked }))
  }

  return (
    <>
      <div className="space-y-4">
        <div className="flex justify-between">
          <Title order={3}>Global Configuration</Title>
          <ActionIcon variant="subtle" disabled={!changed}>
            <LuSave />
          </ActionIcon>
        </div>
        <Card>
          <div className="space-y-8">
            <SettingBoolean
              title="Allow Demo Uploads"
              description="Enable users to manually upload demo files"
              checked={settings.uploads}
              onChange={checked => handleChange("uploads", checked)}
            />
            <SettingBoolean
              title="Custom Highlight Criteria"
              description="Allow regular users to modify their personal highlight generation rules"
              checked={settings.customCriteria}
              onChange={checked => handleChange("customCriteria", checked)}
            />
            <SettingBoolean
              title="Enable Chat Commands"
              description="Listen for specific phrases in game chat to mark rounds for highlights"
              checked={settings.chatCommands}
              onChange={checked => handleChange("chatCommands", checked)}
            />
          </div>
        </Card>
      </div>
      <Block
        shouldBlockFn={() => changed}
        withResolver
      >
        {({ status, proceed, reset }) => (
          <ModalCenter opened={status === "blocked"} onClose={() => reset?.()} size="md">
            <SettingsConfirm proceed={proceed} reset={reset} />
          </ModalCenter>
        )}
      </Block>
    </>
  )
}

const SettingsConfirm = ({ proceed, reset }: { proceed?: () => void, reset?: () => void }) => {
  return (
    <div className="text-center w-full">
      <p>You have unchanged changes</p>
      <p>Are you sure you want to leave?</p>
      <div className="flex justify-end gap-2 mt-8">
        <Button onClick={reset} variant="outline">
          Go Back
        </Button>
        <Button onClick={proceed}>
          Proceed
        </Button>
      </div>
    </div>
  )
}

type SettingBooleanProps = {
  title: string;
  description: string;
  checked: boolean;
  onChange: (value: boolean) => void;
}

const SettingBoolean = ({ title, description, checked, onChange }: SettingBooleanProps) => {
  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-col gap-2">
        <p className="text-white">{title}</p>
        <p className="text-sm text-secondary">{description}</p>
      </div>
      <Switch
        checked={checked}
        onChange={e => onChange(e.target.checked)}
        size="md"
      />
    </div>
  )
}
