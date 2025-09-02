import type { ReactNode } from "react"
import { Header } from "@/components/auth/Header"
import Footer from "@/components/shared/Footer"

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="flex-1">{children}</main>
      {/* mt-auto ensures it sticks to bottom, no extra gap */}
      <footer className="mt-auto">
        <Footer />
      </footer>
    </div>
  )
}
