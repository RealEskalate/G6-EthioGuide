// app/api/auth/callback/google/page.tsx
"use client";

import { Suspense } from "react";
import GoogleCallback from "./GoogleCallbackClient";

export default function GoogleCallbackWrapper() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <GoogleCallback />
    </Suspense>
  );
}
