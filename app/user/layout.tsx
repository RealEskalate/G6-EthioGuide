import type React from "react"
import { Header } from "../../components/shared/Header"
import { Sidebar } from "../../components/shared/Sidebar"
import { Footer } from "../../components/shared/Footer"

export default function UserLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
      <div className="flex flex-col min-h-screen">
      <Header />
      <div className="flex flex-1">
        <Sidebar />
  <main className="flex-1 p-4 md:p-8">{children}</main>
      </div>
      <Footer />
    </div>
  )
}
