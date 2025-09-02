"use client";

import { useEffect } from "react";
import { ChevronDown, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import Image from "next/image";
import { useTranslation } from "react-i18next";

// Language map for display
const languageMap = {
  en: { name: "English", shortCode: "EN" },
  am: { name: "አማርኛ (Amharic)", shortCode: "አማ" },
};

export function Header() {
  const { i18n } = useTranslation();

  // Sync language changes
  const handleLanguageChange = (langCode: string) => {
    i18n.changeLanguage(langCode);
  };

  return (
    <header className="bg-white px-4 py-2 sm:px-6 sm:py-3 sticky top-0 z-50">
      <div className="flex items-center justify-between">
        {/* Logo */}
        <div className="flex items-center gap-2">
          <Image
            src="/images/ethioguide-symbol.png"
            alt="EthioGuide Symbol"
            width={32}
            height={32}
            className="h-8 w-8 sm:h-10 sm:w-10"
            priority
          />
          <span className="text-gray-800 font-semibold text-base sm:text-lg">EthioGuide</span>
        </div>

        {/* Right Section */}
        <div className="flex items-center gap-2 sm:gap-4">
          {/* Language Toggle */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="text-gray-600 hover:text-gray-900 hover:bg-gray-100 sm:text-sm"
              >
                {languageMap[i18n.language as keyof typeof languageMap]?.shortCode || i18n.language}
                <ChevronDown className="w-3 h-3 sm:w-4 sm:h-4 ml-1" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="text-sm sm:text-base p-1 sm:p-2">
              {(Object.keys(languageMap) as Array<keyof typeof languageMap>).map((langCode) => (
                <DropdownMenuItem
                  key={langCode}
                  className="hover:bg-gray-100 hover:text-gray-900 text-xs sm:text-sm"
                  onClick={() => handleLanguageChange(langCode)}
                >
                  {languageMap[langCode].name}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <Button
            className="text-white bg-primary sm:text-sm"
            onClick={() => (window.location.href = "/auth/login")}
          >
            Signin
          </Button>
        </div>
      </div>
    </header>
  );
}