"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Image from "next/image";
import Link from "next/link";
import { MailCheck } from "lucide-react"; 

export default function CheckEmailPage() {
  return (
    <div className="bg-neutral-light text-foreground min-h-[73dvh] flex flex-col flex-1 items-center p-4 sm:pt-6 space-y-2">
      <div className="flex items-center gap-3 p-5">
        <Image
          src="/images/ethioguide-symbol.png"
          alt="EthioGuide Symbol"
          width={50}
          height={50}
          priority
        />
        <span className="text-gray-800 font-semibold text-3xl">EthioGuide</span>
      </div>
      <Card className="bg-background-light w-full max-w-md border-neutral rounded-xl">
        <CardHeader>
          <div className="flex flex-col items-center space-y-4">
            <div className="text-primary rounded-full p-3 bg-primary-light">
              <MailCheck size={48} className="text-primary" /> {/* Using MailCheck icon */}
            </div>
            <CardTitle className="text-2xl font-bold text-center">
              Check Your Email
            </CardTitle>
          </div>
          <p className="text-sm text-center text-neutral-dark mt-2">
            We&#39;ve sent a verification link to your email address.
            Please check your inbox (and spam folder) to activate your account.
          </p>
        </CardHeader>
        <CardContent>
          {/* <p className="mt-4 text-sm text-center text-neutral-dark">
            Didn&#39;t receive the email?{" "}
            <Link href="/auth/resend-verification" className="text-primary hover:underline">
              Resend email
            </Link> */}
          {/* </p> */}
          <p className="mt-2 text-sm text-center text-neutral-dark">
            Already verified?{" "}
            <Link href="/auth/login" className="text-primary hover:underline">
              Login
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}