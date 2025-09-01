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
import { LanguageSwitcher } from "@/components/language-switcher";
import { useSearchParams } from "next/navigation";
import { newPasswordSchema, NewPasswordFormData } from "@/lib/validation/new-password";
import Image from "next/image";
import { FaEye, FaEyeSlash } from "react-icons/fa";
import { useState } from "react";

export default function NewPasswordPage() {
  const { t, i18n } = useTranslation("auth");
  const searchParams = useSearchParams();
  let token = searchParams.get("token"); // Get reset token from URL
  token = "placeholder-token"; // Placeholder token for testing
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const form = useForm<NewPasswordFormData>({
    resolver: zodResolver(newPasswordSchema),
    defaultValues: {
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (data: NewPasswordFormData) => {
    try {
      // Placeholder: Replace with your backend API call
      const response = await fetch("/api/auth/new-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...data, token }),
      });
      if (response.ok) {
        form.setValue("password", "");
        form.setValue("confirmPassword", "");
        form.setError("root", { message: t("new_password.success") });
      } else {
        form.setError("root", { message: t("new_password.error") });
      }
    } catch (error) {
      form.setError("root", { message: t("new_password.error") });
    }
  };

  // // Debug translation loading
  // console.log("Current language:", i18n.language);
  // console.log("Auth translations:", i18n.getResourceBundle(i18n.language, "auth"));

  return (
    <div className="bg-neutral-light text-foreground min-h-screen flex flex-col items-center justify-center p-4 font-sans space-y-4">
      <Image
        src="/logo/logo.png"
        alt="EthioGuide Logo"
        width={240}
        height={240}
        className="object-contain"
        priority
      />
      <Card className="bg-background-light w-full max-w-md border-neutral">
        <CardHeader>
          <div className="flex flex-col items-center space-y-4">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl font-bold text-center">
                {t("new_password.title")}
              </CardTitle>
              {/* <LanguageSwitcher /> */}
            </div>
          </div>
          <p className="text-sm text-center text-neutral-dark ">{t("new_password.sub_title")}</p>
        </CardHeader>
        <CardContent>
          {!token ? (
            <p className="text-error text-center">{t("new_password.invalid_token")}</p>
          ) : (
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic">{t("new_password.password")}</FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showPassword ? "text" : "password"}
                            placeholder={t("new_password.password_placeholder")}
                            className="border-neutral focus:border-primary pr-10"
                            {...field}
                          />
                          <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute inset-y-0 right-0 flex items-center pr-3 text-neutral hover:text-primary"
                            aria-label={showPassword ? "Hide password" : "Show password"}
                          >
                            {showPassword ? <FaEyeSlash className="h-5 w-5" /> : <FaEye className="h-5 w-5" />}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-error" />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic">{t("new_password.confirm_password")}</FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showConfirmPassword ? "text" : "password"}
                            placeholder={t("new_password.confirm_password_placeholder")}
                            className="border-neutral focus:border-primary pr-10"
                            {...field}
                          />
                          <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="absolute inset-y-0 right-0 flex items-center pr-3 text-neutral hover:text-primary"
                            aria-label={showConfirmPassword ? "Hide confirm password" : "Show confirm password"}
                          >
                            {showConfirmPassword ? <FaEyeSlash className="h-5 w-5" /> : <FaEye className="h-5 w-5" />}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-error" />
                    </FormItem>
                  )}
                />
                {form.formState.errors.root && (
                  <p className="text-error text-sm">{form.formState.errors.root.message}</p>
                )}
                <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md"
                >
                  {t("new_password.submit")}
                </Button>
              </form>
            </Form>
          )}
          <p className="mt-2 text-sm text-center text-neutral-dark">
            {t("new_password.back_to_login")}{" "}
            <Link href="/auth/login" className="text-primary hover:underline">
              {t("new_password.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}