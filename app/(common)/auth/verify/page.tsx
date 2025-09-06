"use client";
 
import { Suspense } from "react";
import VerifyComponent from "@/components/common/verifyComponent";

export default function VerifyEmailPage() {
  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <VerifyComponent />
      </Suspense>
    </div>
  );
}