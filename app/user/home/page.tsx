"use client"

import { Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import Image from "next/image"
import Link from "next/link"

export default function UserHomePage() {
  return (
    <div className="min-h-screen w-full bg-gray-50 p-4 ">
  {/* Welcome Section */}
  <div className="mb-8 w-full">
        <h1 className="text-2xl md:text-3xl font-bold text-gray-800 mb-2 text-balance">Welcome to EthioGuide!</h1>
        <p className="text-gray-600 text-pretty">
          Your trusted partner for navigating Ethiopian government procedures with ease.
        </p>
      </div>

  {/* Search Bar */}
  <div className="relative mb-8 w-full" style={{ animationDelay: "0.1s" }}>
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
        <Input
          placeholder="Search government procedures..."
          className="pl-10 py-3 bg-white border-gray-200 text-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
        />
      </div>

  {/* Quick Access Procedures */}
  <section className="mb-8 w-full" style={{ animationDelay: "0.2s" }}>
        <h2 className="text-xl font-semibold text-gray-800 mb-6">Quick Access Procedures</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
          {/* Passport Renewal */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
              <Image src="/icons/business.svg" alt="Passport" width={24} height={24} className="w-6 h-6" />
            </div>
            <h3 className="font-medium text-gray-800 mb-4">Passport Renewal</h3>
            <Link href="/user/procedures-list">
              <Button className="w-full bg-[#6986af] hover:bg-[#3a6a8d] text-white transition-colors duration-200">
                Start Procedure
              </Button>
            </Link>
          </div>

          {/* Business Registration */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
            <div className="w-12 h-12 bg-white rounded-lg flex items-center justify-center mb-4">
              <Image
                src="/icons/business-registration.svg"
                alt="Business Registration"
                width={24}
                height={24}
                className="w-6 h-6"
              />
            </div>
            <h3 className="font-medium text-gray-800 mb-4">Business Registration</h3>
            <Link href="/user/procedures-list">
              <Button className="w-full bg-[#6986af] hover:bg-[#3a6a8d] text-white transition-colors duration-200">
                Start Procedure
              </Button>
            </Link>
          </div>

          {/* National ID Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
              <Image src="/icons/passport.svg" alt="National ID Card" width={24} height={24} className="w-6 h-6" />
            </div>
            <h3 className="font-medium text-gray-800 mb-4">National ID Card</h3>
            <Link href="/user/procedures-list">
              <Button className="w-full bg-[#6986af] hover:bg-[#3a6a8d] text-white transition-colors duration-200">
                Start Procedure
              </Button>
            </Link>
          </div>

          {/* Driver's License */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
            <div className="w-12 h-12 bg-orange-100 rounded-lg flex items-center justify-center mb-4">
              <Image
                src="/icons/driver-license.svg"
                alt="Driver's License"
                width={24}
                height={24}
                className="w-6 h-6"
              />
            </div>
            <h3 className="font-medium text-gray-800 mb-4">Driver&apos;s License</h3>
            <Button className="w-full bg-[#6986af] hover:bg-[#3a6a8d] text-white transition-colors duration-200">
              Start Procedure
            </Button>
          </div>
        </div>
      </section>

  {/* Recent Activity Feed */}
  <section className="w-full" style={{ animationDelay: "0.3s" }}>
        <h2 className="text-xl font-semibold text-gray-800 mb-6">Current Progress</h2>
        <div className="bg-white rounded-lg border border-gray-200 p-6 shadow-sm">
          <div className="space-y-4">
            {/* Activity Items (only completed and progress) */}
            {[
              {
                status: "completed",
                title: "Passport Renewal",
                statusText: "Completed",
                time: "2 hours ago",
                color: "green"
              },
              {
                status: "progress",
                title: "Business License Application",
                statusText: "In Progress",
                time: "1 day ago",
                color: "blue"
              },
              {
                status: "completed",
                title: "Vehicle Registration",
                statusText: "Completed",
                time: "1 week ago",
                color: "green"
              }
            ].map((item, index) => (
              <div
                key={index}
                className="flex items-center justify-between py-3 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 rounded-lg px-2 transition-colors duration-200"
              >
                <div className="flex items-center gap-3">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${
                      item.color === "green"
                        ? "bg-green-100"
                        : "bg-blue-100"
                    }`}
                  >
                    <div
                      className={`w-4 h-4 flex items-center justify-center text-xs font-bold ${
                        item.color === "green"
                          ? "text-green-600"
                          : "text-blue-600"
                      }`}
                    >
                      {item.status === "completed"
                        ? "✓"
                        : "i"}
                    </div>
                  </div>
                  <div>
                    <p className="font-medium text-gray-800">{item.title}</p>
                    <p
                      className={`text-sm ${
                        item.color === "green"
                          ? "text-green-600"
                          : "text-blue-600"
                      }`}
                    >
                      {item.statusText}
                    </p>
                  </div>
                </div>
                <span className="text-sm text-gray-500 whitespace-nowrap">{item.time}</span>
              </div>
            ))}
          </div>

          <div className="mt-6 pt-4 border-t border-gray-100 flex justify-center">
            <Link href="/user/workspace">
              <Button variant="link" className="text-[#3a6a8d] hover:text-[#2d5a7b] transition-colors duration-200">
                View All Activities
              </Button>
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
