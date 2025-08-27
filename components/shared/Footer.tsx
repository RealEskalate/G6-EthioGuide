import { Facebook, Twitter, Instagram, Youtube } from "lucide-react"
import Image from "next/image"

export function Footer() {
  return (
    <footer className="bg-[#2e4d57] text-white py-8 mt-12">
      <div className="max-w-7xl mx-auto px-8">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <Image
              src="/images/ethioguide-symbol.png"
              alt="EthioGuide Symbol"
              width={40}
              height={40}
              className="h-10 w-10"
            />
            <span className="font-semibold text-xl">EthioGuide</span>
          </div>

          {/* Navigation Links */}
          <div className="flex items-center gap-8">
            <a href="#" className="text-[#a7b3b9] hover:text-white text-sm transition-colors duration-200">
              About Us
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white text-sm transition-colors duration-200">
              Contact
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white text-sm transition-colors duration-200">
              Privacy Policy
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white text-sm transition-colors duration-200">
              Terms of Service
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white text-sm transition-colors duration-200">
              Help Center
            </a>
          </div>

          {/* Social Icons */}
          <div className="flex items-center gap-4">
            <a href="#" className="text-[#a7b3b9] hover:text-white transition-colors duration-200">
              <Facebook className="w-5 h-5" />
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white transition-colors duration-200">
              <Instagram className="w-5 h-5" />
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white transition-colors duration-200">
              <Twitter className="w-5 h-5" />
            </a>
            <a href="#" className="text-[#a7b3b9] hover:text-white transition-colors duration-200">
              <Youtube className="w-5 h-5" />
            </a>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-[#4b5563] text-center">
          <p className="text-[#a7b3b9] text-sm">© 2025 EthioGuide. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}
