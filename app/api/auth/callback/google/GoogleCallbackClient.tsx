// app/auth/google/callback/page.tsx
"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { signIn } from "next-auth/react";

export default function GoogleCallback() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const code = searchParams.get("code");
    const errorParam = searchParams.get("error");
    console.log("Google callback params:", { code, error: errorParam });
    if (errorParam) {
      setError("Login failed. Please try again.");
      return;
    }

    if (code) {
      // This function will handle the API call to your backend
      const handleSocialLogin = async () => {
        try {
          // Send the code to your dedicated backend
          const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/auth/social`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ provider: "google", code }),
            }
          );

          if (!res.ok) {
            const errorData = await res.json();
            throw new Error(
              errorData.message || "An error occurred during social login."
            );
          }

          const data = await res.json();
          console.log("Social login data:", data);

          // `data` should be the user object and any tokens from YOUR backend
          // e.g., { user: { id, name, email }, accessToken: 'your_backend_jwt' }

          // Now, sign in to NextAuth using the 'credentials' provider
          // We pass the user data we got from our backend to establish the session.
          const result = await signIn("credentials", {
            redirect: false,
            accessToken: data.access_token,
            role: data.user.role,
            id: data.user.id,
            name: data.user.name,
            email: data.user.email,
          });

          if (result?.error) {
            throw new Error(result.error);
          }

          // Redirect to a protected page after successful login
          router.push("/");
        } catch (err: unknown) {
          console.error("Social login failed:", err);
          if (err instanceof Error) {
            setError(
              err.message || "Failed to process login. Please try again."
            );
          } else {
            setError("Failed to process login. Please try again.");
          }
          // Optional: redirect to login page with an error message
          // router.push('/login?error=SocialLoginFailed');
        }
      };

      handleSocialLogin();
    }
  }, [searchParams, router]);

  // Render a loading or error state
  if (error) {
    return <div style={{ padding: "2rem", color: "red" }}>Error: {error}</div>;
  }

  return <div style={{ padding: "2rem" }}>Processing your login...</div>;
}
