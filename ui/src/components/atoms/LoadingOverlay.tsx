import { cn } from "@/lib/utils";
import { useEffect, useRef, useState } from "react";
import { FragtapeIcon } from "../icons/FragtapeIcon";

type Props = {
  loading: boolean;
};

const SHOW_DELAY_MS = 300;
const MIN_VISIBLE_MS = 150;

export const LoadingOverlay = ({ loading }: Props) => {
  const [visible, setVisible] = useState(false);
  const shownAtRef = useRef<number | null>(null);

  const showTimeoutRef = useRef<number | null>(null);
  const hideTimeoutRef = useRef<number | null>(null);

  useEffect(() => {
    if (loading) {
      if (hideTimeoutRef.current) {
        window.clearTimeout(hideTimeoutRef.current);
        hideTimeoutRef.current = null;
      }

      if (!visible && !showTimeoutRef.current) {
        showTimeoutRef.current = window.setTimeout(() => {
          setVisible(true);
          shownAtRef.current = Date.now();
          showTimeoutRef.current = null;
        }, SHOW_DELAY_MS);
      }
    } else {
      if (showTimeoutRef.current) {
        window.clearTimeout(showTimeoutRef.current);
        showTimeoutRef.current = null;
      }

      if (visible) {
        const shownAt = shownAtRef.current ?? Date.now();
        const elapsed = Date.now() - shownAt;
        const remaining = Math.max(0, MIN_VISIBLE_MS - elapsed);

        hideTimeoutRef.current = window.setTimeout(() => {
          setVisible(false);
          shownAtRef.current = null;
          hideTimeoutRef.current = null;
        }, remaining);
      }
    }

    return () => {
      if (showTimeoutRef.current) {
        window.clearTimeout(showTimeoutRef.current);
        showTimeoutRef.current = null;
      }
      if (hideTimeoutRef.current) {
        window.clearTimeout(hideTimeoutRef.current);
        hideTimeoutRef.current = null;
      }
    };
  }, [loading, visible]);

  return (
    <div
      className={cn(
        "absolute inset-0 z-100 grid place-items-center bg-(--mantine-color-background-7)/60 backdrop-blur-xs origin-left transition-all duration-150",
        visible ? "scale-x-100 opacity-100" : "pointer-events-none scale-x-0 opacity-0"
      )}
    >
      <FragtapeIcon animated className="size-12 text-(--mantine-color-primary-6)" />
    </div>
  );
};
