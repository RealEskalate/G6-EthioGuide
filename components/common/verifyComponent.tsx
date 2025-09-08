"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation"; 
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle2, XCircle, MailWarning } from "lucide-react"; 

export default function VerifyComponent() {
  const searchParams = useSearchParams(); 
  const [verificationStatus, setVerificationStatus] = useState<
    "idle" | "loading" | "success" | "error"
  >("idle");
  const [message, setMessage] = useState("");

  const API = process.env.NEXT_PUBLIC_API_URL;

  useEffect(() => {
    const activateUser = async () => {
      setVerificationStatus("loading");
      const activationToken = searchParams.get("token"); 

      if (!activationToken) {
        setVerificationStatus("error");
        setMessage("No activation token found in the URL.");
        return;
      }

      try {
        const response = await fetch(`${API}/auth/verify`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ activationToken: activationToken }), 
        });

        const data = await response.json();
        console.log("backend",data)
        if (response.ok) {
          setVerificationStatus("success");
          setMessage(data.message || "Email verified successfully!");
        } else {
          setVerificationStatus("error");
          setMessage(
            data.message ||
              "Email verification failed. Please try again or request a new link."
          );
        }
      } catch (error) {
        console.error("Verification API call failed:", error);
        setVerificationStatus("error");
        setMessage("An unexpected error occurred during verification.");
      }
    };

    if (verificationStatus === "idle") {
      activateUser();
    }
  }, [searchParams, verificationStatus]); 
  const getIcon = () => {
    switch (verificationStatus) {
      case "loading":
        return <Loader2 size={48} className="animate-spin text-primary" />;
      case "success":
        return <CheckCircle2 size={48} className="text-green-500" />;
      case "error":
        return <XCircle size={48} className="text-red-500" />;
      default:
        return <MailWarning size={48} className="text-primary" />;
    }
  };

  const getTitle = () => {
    switch (verificationStatus) {
      case "loading":
        return "Verifying Your Email...";
      case "success":
        return "Email Verified!";
      case "error":
        return "Verification Failed";
      default:
        return "Verify Your Email Address";
    }
  };

  return (
    <div className="bg-neutral-light text-foreground min-h-[73dvh] flex flex-col flex-1 items-center p-4 sm:pt-6 space-y-2">
      <div className="flex items-center gap-3 p-5">
        <span className="text-gray-800 font-semibold text-3xl">EthioGuide</span>
      </div>
      <Card className="bg-background-light w-full max-w-md border-neutral rounded-xl">
        <CardHeader>
          <div className="flex flex-col items-center space-y-4">
            <div
              className={`rounded-full p-3 ${
                verificationStatus === "success"
                  ? "bg-green-100"
                  : verificationStatus === "error"
                  ? "bg-red-100"
                  : "bg-primary-light"
              }`}
            >
              {getIcon()}
            </div>
            <CardTitle className="text-2xl font-bold text-center">
              {getTitle()}
            </CardTitle>
          </div>
          <p className="text-sm text-center text-neutral-dark mt-2">
            {verificationStatus === "idle" || verificationStatus === "loading"
              ? "Please wait while we confirm your email address."
              : message}
          </p>
        </CardHeader>
        <CardContent>
          {verificationStatus === "success" && (
            <p className="mt-4 text-sm text-center text-neutral-dark">
              You can now proceed to login.{" "}
              <a href="/auth/login" className="text-primary hover:underline">
                Login
              </a>
            </p>
          )}

          {verificationStatus === "error" && (
            <p className="mt-4 text-sm text-center text-neutral-dark">
              If the problem persists, please contact support.{" "}
              <a href="/auth/login" className="text-primary hover:underline">
                Back to Login
              </a>
            </p>
          )}

          {/* Optional: Add a button to retry if it was an error, though usually with tokens you don't retry the same token */}
          {/* {verificationStatus === "error" && (
            <Button
              onClick={() => {
                setVerificationStatus("idle"); // This will trigger the effect again
                setMessage("");
              }}
              className="w-full mt-4 bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md"
            >
              Retry Verification
            </Button>
          )} */}
        </CardContent>
      </Card>
    </div>
  );
}
