"use client"

import type { ReactNode } from "react"
import React, { Suspense } from "react"
import { Header } from "@/components/auth/Header"
import { Footer } from "@/components/shared/Footer"
// add i18n client init and provider
import i18n from "i18next"
import { I18nextProvider, initReactI18next } from "react-i18next"

// initialize a minimal i18n instance on the client
if (!i18n.isInitialized) {
  i18n.use(initReactI18next).init({
    resources: {},
    lng: "en",
    fallbackLng: "en",
    interpolation: { escapeValue: false },
  })
}

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <I18nextProvider i18n={i18n}>
      <div className="flex flex-col min-h-screen">
        <Suspense fallback={null}>
          <Header />
        </Suspense>
        <main className="flex-1">
          <Suspense fallback={null}>
            {children}
          </Suspense>
        </main>
        <footer className="mt-auto">
          <Footer />
        </footer>
      </div>
    </I18nextProvider>
  )
}
