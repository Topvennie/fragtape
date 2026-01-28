import { Modal, ModalProps } from "@mantine/core";

type Props = ModalProps

export const ModalCenter = (props: Props) => {
  return <Modal
    centered
    size="xl"
    radius="lg"
    {...props}
  />
}
