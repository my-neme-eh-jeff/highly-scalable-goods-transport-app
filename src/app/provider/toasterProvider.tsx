"use client";

import { useTheme } from "next-themes";
import { Toaster } from "sonner";

export function ToasterProvider() {
  const { theme } = useTheme();
  return (
    <Toaster
      richColors
      pauseWhenPageIsHidden
      expand
      theme={theme as "light" | "dark" | "system" | undefined}
    />
  );
}
