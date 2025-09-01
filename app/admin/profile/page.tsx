"use client";
import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { CheckCircle, LogOut, Trash2 } from "lucide-react";
import { FaCrown } from "react-icons/fa";
import rickProfile from "@/public/images/rickProfile.webp";
import PasswordInputBox from "@/components/admin/PasswordInputBox";
import { FaCamera } from "react-icons/fa6";

export default function AccountSettingsPage() {
  const [fullName, setFullName] = useState("Sarah Johnson");
  const [email, setEmail] = useState("sarah.johnson@email.com");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const userInfo = {
    fullName: "Sarah Johnson",
    email: "sarah.johnson@email.com",
  };

  return (
    <div className="min-h-screen bg-gray-50 flex justify-center py-10 px-4 text-primary-dark">
      <div className="w-full max-w-5xl grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left side */}
        <div className="lg:col-span-2 space-y-6">
          {/* Profile Section */}
          <Card className="border-neutral-100">
            <CardContent className="relative flex items-center gap-4 p-6">
              <img
                src={rickProfile.src}
                alt="Profile"
                className="w-20 h-20 rounded-full object-cover"
              />
              {/* change profile picture btn  */}
              <Button className="size-7 absolute bottom-5 left-20 bg-primary rounded-full"><FaCamera className="text-white size-3"/></Button>
              <div>
                <h2 className="text-xl font-semibold">{userInfo.fullName}</h2>
                <p className="text-sm text-gray-500">{userInfo.email}</p>
                <span className="inline-flex items-center 
                p-1 text-secondary px-3 bg-secondary-light rounded-2xl bg-opacity-20 text-sm mt-1">
                  <FaCrown size={16} className="mr-1" /> Paid User
                </span>
              </div>
            </CardContent>
          </Card>

          {/* Account Settings */}
          <Card className="border-neutral-100">
            <CardContent className="p-6 space-y-6">
              <h3 className="text-lg font-semibold">Account Settings</h3>

              {/* Personal Information */}
              <div className="space-y-4">
                <h4 className="font-medium">Personal Informaion</h4>
                <div>
                  <Label htmlFor="fullName" className="text-neutral-dark  ">
                    Full Name
                  </Label>
                  <Input
                    id="fullName"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                  />
                </div>
                <div>
                  <Label htmlFor="email" className="text-neutral-dark  ">
                    Email
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div className="flex gap-3">
                  <Button className="text-white">Save Changes</Button>
                  <Button variant="outline">Cancel</Button>
                </div>
              </div>

              <Separator />

              {/* Change Password */}
              <div className="space-y-4">
                <h4 className="font-medium">Change Password</h4>
                <PasswordInputBox label="Current Password" value={currentPassword} setPass={setCurrentPassword} placeHolder="Enter your current password" />
                <PasswordInputBox label="New Password" value={newPassword} setPass={setNewPassword} placeHolder="Enter your new password" />
                <PasswordInputBox label="Confirm New Password" value={confirmPassword} setPass={setConfirmPassword} placeHolder="Enter your new password" />

                <Button className="bg-primary text-white">
                  Update Password
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right side */}
        <div className="space-y-6">
          {/* Subscription Plan */}
          <Card className="border-neutral-100">
            <CardContent className="p-6">
              <h3 className="font-semibold mb-2">Subscription Plan</h3>
              <div className="relative bg-gradient-to-r from-primary-light to-secondary-light text-white p-4 rounded-xl">
                <h4 className="font-bold">Premium Plan</h4>
                <p>$29/month</p>
                <span><FaCrown className="text-3xl absolute right-5 top-5"/></span>
              </div>
              <ul className="mt-4 text-sm text-gray-600 space-y-2">
                <li className="flex justify-between">
                  All Chat Access:{" "}
                  <span className="font-medium">Unlimited</span>
                </li>
                <li className="flex justify-between">
                  Document Storage: <span className="font-medium">50GB</span>
                </li>
                <li className="flex items-center gap-1 justify-between">
                  Priority Support{" "}
                  <CheckCircle size={16} className="text-green-500" />
                </li>
              </ul>
            </CardContent>
          </Card>

          {/* Logout + Delete */}
          <div className="space-y-3">
            <Button
              variant="outline"
              className="w-full flex items-center justify-center gap-2 border-neutral-200"
            >
              <LogOut size={16} /> Logout
            </Button>
            <Button
              variant="ghost"
              className="w-full text-red-600 hover:text-red-700 flex items-center justify-center gap-2"
            >
              <Trash2 size={16} /> Delete Account
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
