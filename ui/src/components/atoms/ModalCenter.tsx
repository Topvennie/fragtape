import { cn } from "@/lib/utils";
import { Modal, ModalProps } from "@mantine/core";

type Props = ModalProps

export const ModalCenter = ({ className, ...props }: Props) => {
  return <Modal
    centered
    size="xl"
    radius="lg"
    className={cn("text-white", className)}
    {...props}
  />
}
