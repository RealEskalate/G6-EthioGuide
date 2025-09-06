"use client";

import { useEffect, useState } from "react";
import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { NotificationDropdown } from "./NotificationDropdown";
import Image from "next/image";
import { usePathname } from "next/navigation";

export function Header() {
  const [language, setLanguage] = useState("EN");
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  const pathname = usePathname();
  const isAdmin = pathname.startsWith(`/admin`);
  const isOrg = pathname.startsWith(`/organization`);
  return (
    <header className="bg-white px-6 py-4 sticky top-0 z-50">
      <div className="flex items-center justify-between">
        {/* Logo */}
        <div className="flex items-center gap-3">
          <Image
            src="/images/ethioguide-symbol.png"
            alt="EthioGuide Symbol"
            width={40}
            height={40}
            className="h-10 w-10"
            priority
          />
          <span className="text-gray-800 font-semibold text-xl">
            EthioGuide
          </span>
        </div>

        {/* Right Section */}
        <div className="flex items-center gap-4">
          {/* Language Toggle */}
          {mounted ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-gray-600 hover:text-gray-900 hover:bg-gray-100"
                >
                  {language}
                  <ChevronDown className="w-4 h-4 ml-1" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  className="hover:bg-gray-100 hover:text-gray-900"
                  onClick={() => setLanguage("EN")}
                >
                  English
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="hover:bg-gray-100 hover:text-gray-900"
                  onClick={() => setLanguage("አማ")}
                >
                  አማርኛ (Amharic)
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Button variant="ghost" size="sm" className="text-gray-600" aria-hidden>
              {language}
              <ChevronDown className="w-4 h-4 ml-1" />
            </Button>
          )}

          {/* Notifications */}
          {mounted ? <NotificationDropdown /> : (
            <Button variant="ghost" size="sm" className="relative p-2 rounded-full" aria-hidden>
              <Image src="/icons/notifications.svg" alt="Notifications" width={20} height={20} className="w-5 h-5" />
            </Button>
          )}

          {/* Profile */}
          <Button
            variant="ghost"
            size="sm"
            className="p-0 rounded-full hover:bg-gray-100"
            onClick={() => {
              if (isAdmin) {
                window.location.href = "/admin/profile";
              }
              else if (isOrg) {
                window.location.href = "/organization/profile";
              }else {
                window.location.href = "/user/profile";
              }
            }}
            aria-label="Go to profile"
          >
            <div className="w-8 h-8 rounded-full overflow-hidden border-2 border-gray-200">
              <Image
                src="/images/profile-photo.jpg"
                alt="Profile Photo"
                width={32}
                height={32}
                className="w-full h-full object-cover"
              />
            </div>
          </Button>
        </div>
      </div>
    </header>
  );
}
