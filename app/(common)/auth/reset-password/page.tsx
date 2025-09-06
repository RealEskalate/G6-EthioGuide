"use client";

import React, { Suspense } from "react";
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
import { useTranslation, initReactI18next } from "react-i18next";
import i18n from "i18next";
import Link from "next/link";
import {
  resetPasswordSchema,
  ResetPasswordFormData,
} from "@/lib/validation/reset-password";
import Image from "next/image";

// init a minimal i18n instance to avoid react-i18next warning
if (!i18n.isInitialized) {
  i18n.use(initReactI18next).init({
    resources: {},
    lng: "en",
    fallbackLng: "en",
    interpolation: { escapeValue: false },
  });
}

function ResetPasswordContent() {
  const { t } = useTranslation("auth");
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

  return (
    <div className="bg-neutral-light text-foreground min-h-[73dvh] flex flex-col flex-1 items-center p-4 sm:pt-6 space-y-2">
      <div className="flex items-center gap-3 p-5">
        <Image
          src="/images/ethioguide-symbol.png"
          alt="EthioGuide Symbol"
          width={50}
          height={50}
          priority
          style={{ width: "auto", height: "auto" }}
        />
        <span className="text-gray-800 font-semibold text-3xl">EthioGuide</span>
      </div>
      <Card className="bg-background-light w-full max-w-md border-neutral">
        <CardHeader>
          <div className="flex flex-col items-center space-y-4">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl font-bold text-center">
                {t("reset_password.title")}
              </CardTitle>
              {/* <LanguageSwitcher /> */}
            </div>
          </div>
          <p className="text-sm text-center text-neutral-dark ">
            {t("reset_password.sub_title")}
          </p>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">
                      {t("reset_password.email")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t("reset_password.email_placeholder")}
                        className="border-neutral focus:border-primary"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-error" />
                  </FormItem>
                )}
              />
              {form.formState.errors.root && (
                <p className="text-error text-sm">
                  {form.formState.errors.root.message}
                </p>
              )}
              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md"
              >
                {t("reset_password.submit")}
              </Button>
            </form>
          </Form>
          <p className="mt-2 text-sm text-center text-neutral-dark">
            {t("reset_password.back_to_login")}{" "}
            <Link href="/auth/login" className="text-primary hover:underline">
              {t("reset_password.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={<div className="p-4 text-gray-600">Loading...</div>}>
      <ResetPasswordContent />
    </Suspense>
  );
}


