import { cn } from "@/lib/utils";
import { Modal as MModal, ModalProps } from "@mantine/core";

type Props = ModalProps

export const Modal = ({ className, ...props }: Props) => {
  return <MModal
    size="xl"
    radius="lg"
    className={cn("text-white", className)}
    {...props}
  />
}

export const ModalCenter = (props: Props) => {
  return <Modal
    centered
    {...props}
  />
}
