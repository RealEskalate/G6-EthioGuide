import { AdminSidebar } from "@/components/shared/AdminSidebar";
import { Header } from "@/components/shared/Header";
import { Footer } from "@/components/shared/Footer";
import { Children, ReactNode } from "react";

export default function AdminLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <div className="bg-gray-50 flex flex-1">
        <AdminSidebar />
        <main className="flex-1 min-w-0">{children}</main>
      </div>
      <Footer />
    </div>
  );
}
