import { useSettingUserFirstTimeWizard, useSettingUserGet } from "@/lib/api/setting_user"
import { useDisclosure } from "@mantine/hooks"
import { ModalCenter } from "../atoms/Modal"
import { Title } from "../atoms/Title"
import { Alert } from "../atoms/Alert"
import { LuArrowRight, LuInfo } from "react-icons/lu"
import { Button, Stack } from "@mantine/core"
import { useDemoGetFiltered } from "@/lib/api/demo"
import { useNavigate } from "@tanstack/react-router"

export const UserFirstTime = () => {
  const { data: setting } = useSettingUserGet()
  const { result: demos } = useDemoGetFiltered()

  const settingUserFirsTimeWizard = useSettingUserFirstTimeWizard()

  const [opened] = useDisclosure(true)

  const navigate = useNavigate()

  if (!setting?.firstTimeWizard) return null

  const handleClose = () => {
    setting.firstTimeWizard = false

    settingUserFirsTimeWizard.mutateAsync(setting)
  }

  const handleSettings = () => {
    handleClose()
    navigate({ to: "/setting" })
  }

  return (
    <ModalCenter opened={opened} onClose={handleClose} withCloseButton={false}>
      <Stack>
        <Title order={2}>
          Welcome to Fragtape!
        </Title>
        <p className="text-secondary">{`You're now ready to start automatically generating highlights from your matches.`}</p>
        {demos.total > 0 && (
          <Alert icon={<LuInfo className="size-6 text-(--mantine-color-primary-6)" />} color="blue">
            <p className="text-white">Why do I already see matches?</p>
            <p className="text-secondary">If another Fragtape user played in one of your matches (as a teammate or opponent), that match is already in our system and we automatically link it to your profile.</p>
          </Alert>
        )}
        <p className="text-secondary">Head over to your settings to manage your account integrations!</p>
        <Button onClick={handleSettings} rightSection={<LuArrowRight className="size-4" />} className="ml-auto">
          {`Got it, let's go`}
        </Button>
      </Stack>
    </ModalCenter>
  )
}
