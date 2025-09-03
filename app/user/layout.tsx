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
      <div className="flex flex-1 flex-col md:flex-row">
        {/* Sidebar: hidden on mobile, visible on md+ screens */}
        <aside className="hidden md:flex md:flex-col flex-shrink-0 h-full">
          <Sidebar className="flex-shrink-0" />
        </aside>
        {/* Mobile sidebar: visible only on mobile, collapsible */}
        <aside className="md:hidden w-full flex flex-col">
          <Sidebar />
        </aside>
        <main className="flex-1 p-2 sm:p-4 md:p-8">{children}</main>
      </div>
      <Footer />
    </div>
  )
}
         
