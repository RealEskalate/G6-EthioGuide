"use client";
import { Suspense } from "react";
import NewPasswordComponent from "@/components/common/newPasswordComponent";
 

export default function NewPasswordPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <NewPasswordComponent />
    </Suspense>
  );
}