// app/organizations/new/page.tsx
"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { FaInfoCircle } from "react-icons/fa";
import { FaGlobe } from "react-icons/fa";
import { BiSolidContact } from "react-icons/bi";
import { IoSettingsSharp } from "react-icons/io5";
import { FaTwitter } from "react-icons/fa";
import { FaInstagram } from "react-icons/fa";
import { FaFacebook } from "react-icons/fa";
import { FaLinkedin } from "react-icons/fa";


export default function CreateOrganizationPage() {
  const [status, setStatus] = useState("active")

  return (
    <div className="p-8 max-w-4xl mx-auto space-y-6 text-primary-dark">
      {/* Page Title */}
      <div>
        <h1 className="text-primary-dark text-2xl font-bold tracking-tight">Create New Organization</h1>
        <p className="text-muted-foreground text-neutral">
          Set up a new organization profile to manage procedures, notices, and feedback.
        </p>
      </div>

      {/* Basic Information */}
      <Card className="border-neutral-100">
        <CardHeader>
          <CardTitle className="flex"> <FaInfoCircle className="text-primary-light mr-2"/>Basic Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label className="p-1">Organization Name *</Label>
            <Input placeholder="Enter organization name"/>
          </div>
          <div>
            <Label className="p-1">Category / Type *</Label>
            <Select>
              <SelectTrigger>
                <SelectValue placeholder="Select category..." />
              </SelectTrigger>
              <SelectContent className="bg-white">
                <SelectItem value="ngo">NGO</SelectItem>
                <SelectItem value="company">Company</SelectItem>
                <SelectItem value="school">School</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label className="p-1">Short Description</Label>
            <Textarea placeholder="Brief description of the organization..." />
          </div>
        </CardContent>
      </Card>

      {/* Contact Information */}
      <Card className="border-neutral-100">
        <CardHeader>
          <CardTitle className="flex"><BiSolidContact className="text-primary-light mr-2"/>Contact Information</CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label className="p-1">Email Address *</Label>
            <Input type="email" placeholder="contact@organization.com" />
          </div>
          <div>
            <Label className="p-1">Phone Number</Label>
            <Input type="tel" placeholder="+1 (555) 123-4567" />
          </div>
          <div>
            <Label className="p-1">Website URL</Label>
            <Input placeholder="https://www.organization.com" />
          </div>
          <div>
            <Label className="p-1">City</Label>
            <Input placeholder="City" />
          </div>
          <div className="md:col-span-2">
            <Label className="p-1">Physical Address</Label>
            <Input placeholder="Street address, city, state, postal code" />
          </div>
        </CardContent>
      </Card>

      {/* Social & Online Presence */}
      <Card className="border-neutral-100">
        <CardHeader>
          <CardTitle className="flex items-end"><FaGlobe className="mr-2 text-primary-light"/>Social & Online Presence <span className="text-sm text-muted-foreground ml-2">(Optional)</span></CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label className="p-1">Organization Tagline</Label>
            <Input placeholder="Your organization's motto or tagline" />
          </div>
          <div>
            <Label className="p-1"><span><FaFacebook className="text-neutral"/></span>Facebook</Label>
            <Input placeholder="Facebook page URL" />
          </div>
          <div>
            <Label className="p-1"><span><FaTwitter className="text-neutral"/></span>Twitter</Label>
            <Input placeholder="Twitter profile URL" />
          </div>
          <div>
            <Label className="p-1"><span><FaLinkedin className="text-neutral"/></span>LinkedIn</Label>
            <Input placeholder="LinkedIn company URL" />
          </div>
          <div className="md:col-span-2">
            <Label className="p-1"><span><FaInstagram className="text-neutral"/></span>Instagram</Label>
            <Input placeholder="Instagram profile URL" />
          </div>
        </CardContent>
      </Card>

      {/* Administration Settings */}
      <Card className="border-neutral-100">
        <CardHeader>
          <CardTitle className="flex"><IoSettingsSharp className="text-primary-light mr-2" />Administration Settings</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label className="p-1">Assign Admin(s) *</Label>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Select admin..." />
                </SelectTrigger>
                <SelectContent className="bg-white">
                  <SelectItem value="admin1">Admin 1</SelectItem>
                  <SelectItem value="admin2">Admin 2</SelectItem>
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground mt-1">
                You can assign multiple administrators
              </p>
            </div>
            <div>
              <Label className="p-1">Default Visibility *</Label>
              <Select defaultValue="public">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent className="bg-white">
                  <SelectItem value="public">Public</SelectItem>
                  <SelectItem value="private">Private</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div>
            <Label className="p-1">Status *</Label>
            <RadioGroup value={status} onValueChange={setStatus} className="flex space-x-4 mt-2">
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="active" id="active" />
                <Label htmlFor="active" className="p-1">Active</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="inactive" id="inactive" />
                <Label htmlFor="inactive" className="p-1">Inactive</Label>
              </div>
            </RadioGroup>
          </div>
        </CardContent>
      </Card>

      {/* Action Buttons */}
      <div className="flex justify-end gap-3">
        <Button variant="outline">Cancel</Button>
        <Button className="bg-primary text-white">Create Organization</Button>
      </div>
    </div>
  )
}
