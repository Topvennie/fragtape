import { LuTriangleAlert, LuArrowRight } from "react-icons/lu"
import { LinkButton } from "../atoms/LinkButton"
import { Alert } from "../atoms/Alert"

export const SettingNoConnections = () => {
  return (
    <Alert
      title="Missing Account Connections"
      icon={<LuTriangleAlert className="size-6 text-(--mantine-color-primary-6)" />}
      color="orange"
      border
    >
      {`Your profile has no connections.\nWe cannot fetch your matches or generate highlights until you add at least one connection.`}
      <div className="mt-4">
        <LinkButton to="/setting" rightSection={<LuArrowRight />}>
          Go to Settings
        </LinkButton>
      </div>
    </Alert>
  )
}
