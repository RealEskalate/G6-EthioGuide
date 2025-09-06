"use client";

import { SessionProvider } from "next-auth/react";
import { I18nextProvider } from "react-i18next";
import i18next from "@/lib/i18n/i18n";
import { ReactNode } from "react";
import { Provider } from "react-redux";
import { store } from "../app/store/store";

interface ProvidersProps {
  children: ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  return (
    <SessionProvider>
      <I18nextProvider i18n={i18next}>
        <Provider store={store}>{children}</Provider>
      </I18nextProvider>
    </SessionProvider>
  );
}
