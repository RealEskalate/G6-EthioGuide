"use client";

import { Card, CardHeader, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useState } from "react";

export default function OrganizationSettings() {
  // Example state, replace with actual data/fetch logic
  const [orgName, setOrgName] = useState("Tech Solutions Inc.");
  const [contactPerson, setContactPerson] = useState("John Smith");
  const [email, setEmail] = useState("contact@techsolutions.com");
  const [city, setCity] = useState("San Francisco");
  const [subCity, setSubCity] = useState("Downtown");
  const [woreda, setWoreda] = useState("District 1");
  const [fullAddress, setFullAddress] = useState("123 Tech Street, Suite 456, San Francisco, CA 94105");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");

  return (
    <div className="min-h-screen bg-[#f7fafd] flex flex-row">
      {/* Sidebar */}
      <aside className="w-[300px] bg-white rounded-xl shadow p-8 flex flex-col items-center mt-8 ml-8 h-fit">
        <div className="flex flex-col items-center mb-8">
          <div className="bg-[#eaf2f7] rounded-full p-4 mb-2">
            {/* Icon placeholder */}
            <svg width="32" height="32" fill="none" viewBox="0 0 24 24">
              <path fill="#2d5c7f" d="M12 2C7.03 2 3 6.03 3 11c0 4.97 4.03 9 9 9s9-4.03 9-9c0-4.97-4.03-9-9-9Zm0 16c-3.87 0-7-3.13-7-7 0-3.87 3.13-7 7-7s7 3.13 7 7c0 3.87-3.13 7-7 7Z"/>
              <circle cx="12" cy="11" r="3" fill="#2d5c7f"/>
            </svg>
          </div>
          <div className="font-semibold text-lg text-[#2d5c7f]">{orgName}</div>
          <Badge variant="secondary" className="mt-1 text-xs bg-[#eaf2f7] text-[#2d5c7f] px-3 py-1 rounded">{`Organization ID: ORG-2024-001`}</Badge>
        </div>
        <nav className="w-full">
          <ul>
            <li className="mb-2">
              <Button variant="outline" className="w-full justify-start text-neutral border-[#eaf2f7] bg-[#eaf2f7] font-semibold">
                <svg className="inline mr-2" width="18" height="18" fill="none" viewBox="0 0 24 24">
                  <path fill="#2d5c7f" d="M12 2C7.03 2 3 6.03 3 11c0 4.97 4.03 9 9 9s9-4.03 9-9c0-4.97-4.03-9-9-9Zm0 16c-3.87 0-7-3.13-7-7 0-3.87 3.13-7 7-7s7 3.13 7 7c0 3.87-3.13 7-7 7Z"/>
                </svg>
                Organization Details
              </Button>
            </li>
            <li>
              <Button variant="ghost" className="w-full justify-start text-neutral font-semibold">
                <svg className="inline mr-2" width="18" height="18" fill="none" viewBox="0 0 24 24">
                  <path fill="#2d5c7f" d="M12 2a10 10 0 100 20 10 10 0 000-20zm0 18a8 8 0 110-16 8 8 0 010 16zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
                </svg>
                Account Settings
              </Button>
            </li>
          </ul>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col gap-8 mt-8 ml-8 mr-8 max-w-4xl">
        {/* Organization Details */}
        <Card className="mb-4 rounded-xl border border-gray-100 shadow-none bg-white">
          <CardHeader className="pb-2 border-b border-gray-100">
            <div className="flex items-center gap-2">
              <svg width="22" height="22" fill="none" viewBox="0 0 24 24">
                <path fill="#2d5c7f" d="M12 2C7.03 2 3 6.03 3 11c0 4.97 4.03 9 9 9s9-4.03 9-9c0-4.97-4.03-9-9-9Zm0 16c-3.87 0-7-3.13-7-7 0-3.87 3.13-7 7-7s7 3.13 7 7c0 3.87-3.13 7-7 7Z"/>
              </svg>
              <h2 className="font-semibold text-lg text-[#2d5c7f]">Organization Details</h2>
            </div>
          </CardHeader>
          <CardContent>
            <form className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-semibold text-neutral mb-2">Organization Name</label>
                  <Input
                    value={orgName}
                    onChange={e => setOrgName(e.target.value)}
                    className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-neutral mb-2">Contact Person</label>
                  <Input
                    value={contactPerson}
                    onChange={e => setContactPerson(e.target.value)}
                    className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-semibold text-neutral mb-2">Email Address</label>
                <Input
                  type="email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                />
              </div>
              <div>
                <div className="font-semibold text-[#2d5c7f] mb-2 mt-4 text-lg">Address Information</div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-4">
                  <div>
                    <label className="block text-sm font-semibold text-neutral mb-2">City</label>
                    <Input
                      value={city}
                      onChange={e => setCity(e.target.value)}
                      className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-neutral mb-2">Sub-city</label>
                    <Input
                      value={subCity}
                      onChange={e => setSubCity(e.target.value)}
                      className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-neutral mb-2">Woreda</label>
                    <Input
                      value={woreda}
                      onChange={e => setWoreda(e.target.value)}
                      className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-semibold text-neutral mb-2">Full Address</label>
                  <textarea
                    value={fullAddress}
                    onChange={e => setFullAddress(e.target.value)}
                    rows={2}
                    className="w-full rounded-lg border border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] p-3 resize-none text-base"
                  />
                </div>
              </div>
              <div className="flex gap-3 mt-4">
                <Button type="submit" className="bg-[#2d5c7f] hover:bg-[#20435c] text-white px-6 py-2 rounded-lg h-12">Save Changes</Button>
                <Button variant="outline" type="button" className="px-6 py-2 rounded-lg h-12 border-gray-300">Cancel</Button>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Account Settings */}
        <Card className="mb-4 rounded-xl border border-gray-100 shadow-none bg-white">
          <CardHeader className="pb-2 border-b border-gray-100">
            <div className="flex items-center gap-2">
              <svg width="22" height="22" fill="none" viewBox="0 0 24 24">
                <path fill="#2d5c7f" d="M12 2a10 10 0 100 20 10 10 0 000-20zm0 18a8 8 0 110-16 8 8 0 010 16zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
              </svg>
              <h2 className="font-semibold text-lg text-[#2d5c7f]">Account Settings</h2>
            </div>
          </CardHeader>
          <CardContent>
            <div className="font-semibold text-[#2d5c7f] mb-2 mt-2 text-base">Update Credentials</div>
            <form className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-semibold text-neutral mb-2">Current Password</label>
                  <Input
                    type="password"
                    value={currentPassword}
                    onChange={e => setCurrentPassword(e.target.value)}
                    className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-neutral mb-2">New Password</label>
                  <Input
                    type="password"
                    value={newPassword}
                    onChange={e => setNewPassword(e.target.value)}
                    className="rounded-lg border-gray-300 focus:border-[#2d5c7f] focus:ring-[#2d5c7f] text-base h-12"
                  />
                </div>
              </div>
              <Button type="submit" className="bg-[#2d5c7f] hover:bg-[#20435c] text-white px-6 py-2 rounded-lg h-12 mt-2">Update Password</Button>
            </form>
          </CardContent>
        </Card>

        {/* Danger Zone */}
        <div className="rounded-xl border border-red-300 bg-white mt-4 p-0 overflow-hidden">
          <div className="flex items-center gap-2 px-6 py-4 border-b border-red-300 bg-red-50">
            <svg width="22" height="22" fill="none" viewBox="0 0 24 24">
              <path fill="#e53e3e" d="M12 2a10 10 0 100 20 10 10 0 000-20zm1 14h-2v-2h2zm0-4h-2V7h2z"/>
            </svg>
            <span className="font-semibold text-lg text-[#e53e3e]">Danger Zone</span>
          </div>
          <div className="flex items-center justify-between px-6 py-4 bg-red-50">
            <div>
              <div className="font-semibold text-[#e53e3e]">Delete Account</div>
              <div className="text-sm text-[#e53e3e]">Permanently delete your organization and all data</div>
            </div>
            <Button variant="destructive" className="bg-[#e53e3e] hover:bg-[#c53030] text-white px-6 py-2 rounded-lg h-12">Delete Account</Button>
          </div>
        </div>
      </main>
    </div>
  );
}
