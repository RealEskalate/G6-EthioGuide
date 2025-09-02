import { Facebook, Twitter, Instagram, Youtube } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

export function Footer() {
  return (
    <footer className="bg-primary-dark text-white py-4 sm:py-6">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          {/* Logo */}
          <div className="flex items-center gap-2 justify-center sm:justify-start">
            <Image
              src="/images/ethioguide-symbol.png"
              alt="EthioGuide Symbol"
              width={32}
              height={32}
              className="h-8 w-8 sm:h-10 sm:w-10"
            />
            <span className="font-semibold text-lg sm:text-xl">EthioGuide</span>
          </div>

          {/* Navigation Links */}
          <div className="flex flex-wrap justify-center gap-4 sm:gap-6">
            <Link
              href="/about"
              className="text-[#a7b3b9] hover:text-primary text-xs sm:text-sm transition-colors duration-200"
            >
              About Us
            </Link>
            <Link
              href="/contact"
              className="text-[#a7b3b9] hover:text-primary text-xs sm:text-sm transition-colors duration-200"
            >
              Contact
            </Link>
            <Link
              href="/privacy"
              className="text-[#a7b3b9] hover:text-primary text-xs sm:text-sm transition-colors duration-200"
            >
              Privacy Policy
            </Link>
            <Link
              href="/terms"
              className="text-[#a7b3b9] hover:text-primary text-xs sm:text-sm transition-colors duration-200"
            >
              Terms of Service
            </Link>
            <Link
              href="/help"
              className="text-[#a7b3b9] hover:text-primary text-xs sm:text-sm transition-colors duration-200"
            >
              Help Center
            </Link>
          </div>

          {/* Social Icons */}
          <div className="flex items-center justify-center gap-3 sm:gap-4">
            <a
              href="https://facebook.com"
              className="text-[#a7b3b9] hover:text-primary transition-colors duration-200"
              aria-label="Follow us on Facebook"
            >
              <Facebook className="w-4 h-4 sm:w-5 sm:h-5" />
            </a>
            <a
              href="https://instagram.com"
              className="text-[#a7b3b9] hover:text-primary transition-colors duration-200"
              aria-label="Follow us on Instagram"
            >
              <Instagram className="w-4 h-4 sm:w-5 sm:h-5" />
            </a>
            <a
              href="https://twitter.com"
              className="text-[#a7b3b9] hover:text-primary transition-colors duration-200"
              aria-label="Follow us on Twitter"
            >
              <Twitter className="w-4 h-4 sm:w-5 sm:h-5" />
            </a>
            <a
              href="https://youtube.com"
              className="text-[#a7b3b9] hover:text-primary transition-colors duration-200"
              aria-label="Follow us on YouTube"
            >
              <Youtube className="w-4 h-4 sm:w-5 sm:h-5" />
            </a>
          </div>
        </div>

        <div className="mt-4 pt-4 border-t border-[#4b5563] text-center">
          <p className="text-[#a7b3b9] text-xs sm:text-sm">
            © {new Date().getFullYear()} EthioGuide. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}