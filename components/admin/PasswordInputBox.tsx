import React from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { EyeOff, Eye } from "lucide-react";

interface passwordProps {
  label: string;
  value: string;
  placeHolder: string;
  setPass: (value: string) => void;
}

const PasswordInputBox = ({ label, value, setPass, placeHolder }: passwordProps) => {
  const [showCurPass, setshowCurPass] = useState(false);

  return (
    <div>
      <Label htmlFor="curPass" className="text-neutral-dark mb-2">
        {label}
      </Label>
      <div className="relative">
        <Input
          id="curPass"
          type={showCurPass ? "text" : "password"}
          value={value}
          onChange={(e) => setPass(e.target.value)}
          className="pr-10"
          placeholder={placeHolder}
        />
        <button
          type="button"
          onClick={() => setshowCurPass(!showCurPass)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-neutral"
        >
          {showCurPass ? (
            <EyeOff className="hover:text-primary-dark" />
          ) : (
            <Eye className="hover:text-primary-dark" />
          )}
        </button>
      </div>
    </div>
  );
};

export default PasswordInputBox;
