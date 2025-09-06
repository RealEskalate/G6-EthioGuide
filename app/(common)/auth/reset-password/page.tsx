"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useTranslation } from "react-i18next";
import Link from "next/link";
import {
  resetPasswordSchema,
  ResetPasswordFormData,
} from "@/lib/validation/reset-password";
import Image from "next/image";

export default function ResetPasswordPage() {
  const { t, i18n } = useTranslation("auth");
  const API_URL = process.env.NEXT_PUBLIC_API_URL;
  const form = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: ResetPasswordFormData) => {
    try {
      // Placeholder: Replace with your backend API call
      const response = await fetch(`${API_URL}/auth/forgot`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      if (response.ok) {
        form.setValue("email", "");
        form.setError("root", { message: t("reset_password.success") });
      } else {
        form.setError("root", { message: t("reset_password.error") });
      }
    } catch {
      form.setError("root", { message: t("new_password.error") });
    }
  };

  // Debug translation loading
  console.log("Current language:", i18n.language);
  console.log(
    "Auth translations:",
    i18n.getResourceBundle(i18n.language, "auth")
  );

  return (
    <div className=" bg-neutral-light text-foreground flex flex-col items-center justify-center px-4 py-6">
      <div className="flex items-center justify-center gap-3 mb-2">
        <Image
          src="/images/ethioguide-symbol.png"
          alt="EthioGuide Symbol"
          width={70}
          height={70}
          priority
        />
        <span className="text-gray-900 font-semibold text-2xl">EthioGuide</span>
      </div>
      <Card className="w-full max-w-md bg-white border border-gray-200 shadow-lg rounded-2xl overflow-hidden">
        <CardHeader className="pb-4">
          <div className="flex flex-col items-center space-y-4">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl sm:text-2xl font-bold text-center font-amharic bg-gradient-to-r from-[#2e4d57] to-[#1c3b2e] bg-clip-text text-transparent">
                {t("reset_password.title")}
              </CardTitle>
            </div>
          </div>
          <p className="text-sm text-center text-[#2e4d57]/80 ">
            {t("reset_password.sub_title")}
          </p>
        </CardHeader>
        <CardContent className="px-6 pb-6">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("reset_password.email")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t("reset_password.email_placeholder")}
                        className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-red-500 text-xs" />
                  </FormItem>
                )}
              />
              {form.formState.errors.root && (
                <p className="text-green-500 text-sm bg-green-50 p-3 rounded-lg border border-green-200">
                  {form.formState.errors.root.message}
                </p>
              )}
              <Button
                type="submit"
                className="w-full h-12 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 hover:shadow-lg focus:ring-4 focus:ring-[#3a6a8d]/30"
                style={{
                  background: `linear-gradient(135deg, #3a6a8d 0%, #2e4d57 50%, #1c3b2e 100%)`,
                }}
              >
                {t("reset_password.submit")}
              </Button>
            </form>
          </Form>
          <p className="mt-6 text-sm text-center text-[#2e4d57]/80">
            {t("reset_password.back_to_login")}{" "}
            <Link
              href="/auth/login"
              className="text-[#3a6a8d] hover:text-[#2e4d57] font-semibold transition-colors duration-200 hover:underline"
            >
              {t("reset_password.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
