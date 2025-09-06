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
import { useSearchParams } from "next/navigation";
import {
  newPasswordSchema,
  NewPasswordFormData,
} from "@/lib/validation/new-password";
import Image from "next/image";
import { FaEye, FaEyeSlash } from "react-icons/fa";
import { useState, Suspense } from "react";

export default function NewPasswordPage() {
  return (
    <Suspense>
      <NewPasswordPageContent />
    </Suspense>
  );
}

function NewPasswordPageContent() {
  const { t } = useTranslation("auth");
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
    } catch {
      form.setError("root", { message: t("new_password.error") });
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 text-foreground flex flex-col items-center justify-center px-4 py-8">
      <Card className="w-full max-w-md bg-white border border-gray-200 shadow-lg rounded-2xl overflow-hidden">
        <CardHeader className="pb-4">
          <div className="flex items-center justify-center gap-3 mb-2">
            <Image
              src="/images/ethioguide-symbol.png"
              alt="EthioGuide Symbol"
              width={40}
              height={40}
              priority
            />
            <span className="text-gray-900 font-semibold text-2xl">EthioGuide</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl sm:text-3xl font-bold text-center font-amharic bg-gradient-to-r from-[#2e4d57] to-[#1c3b2e] bg-clip-text text-transparent">
                {t("new_password.title")}
              </CardTitle>
            </div>
          </div>
          <p className="text-sm text-center text-[#2e4d57]/80 font-medium">
            {t("new_password.sub_title")}
          </p>
        </CardHeader>
        <CardContent className="px-6 pb-6">
          {!token ? (
            <p className="text-red-500 text-sm text-center">
              {t("new_password.invalid_token")}
            </p>
          ) : (
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                        {t("new_password.password")}
                      </FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showPassword ? "text" : "password"}
                            placeholder={t("new_password.password_placeholder")}
                            className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90 pr-12"
                            {...field}
                            onBlur={() => form.trigger("password")}
                          />
                          <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute inset-y-0 right-0 flex items-center pr-4 text-[#2e4d57] hover:text-[#3a6a8d] transition-colors duration-200"
                            aria-label={
                              showPassword ? "Hide password" : "Show password"
                            }
                          >
                            {showPassword ? <FaEyeSlash className="h-5 w-5" /> : <FaEye className="h-5 w-5" />}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-red-500 text-xs" />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                        {t("new_password.confirm_password")}
                      </FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showConfirmPassword ? "text" : "password"}
                            placeholder={t("new_password.confirm_password_placeholder")}
                            className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90 pr-12"
                            {...field}
                            onBlur={() => form.trigger("confirmPassword")}
                          />
                          <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="absolute inset-y-0 right-0 flex items-center pr-4 text-[#2e4d57] hover:text-[#3a6a8d] transition-colors duration-200"
                            aria-label={
                              showConfirmPassword
                                ? "Hide confirm password"
                                : "Show confirm password"
                            }
                          >
                            {showConfirmPassword ? <FaEyeSlash className="h-5 w-5" /> : <FaEye className="h-5 w-5" />}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-red-500 text-xs" />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="w-full h-12 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 hover:shadow-lg focus:ring-4 focus:ring-[#3a6a8d]/30"
                  style={{ background: `linear-gradient(135deg, #3a6a8d 0%, #2e4d57 50%, #1c3b2e 100%)` }}
                >
                  {t("new_password.submit")}
                </Button>
                {form.formState.errors.root && (
                  <p className="text-red-500 text-sm text-center mt-2 bg-red-50 p-3 rounded-lg border border-red-200">
                    {form.formState.errors.root.message}
                  </p>
                )}
              </form>
            </Form>
          )}
          <p className="mt-6 text-sm text-center text-[#2e4d57]/80">
            {t("new_password.back_to_login")}{" "}
            <Link href="/auth/login" className="text-[#3a6a8d] hover:text-[#2e4d57] font-semibold transition-colors duration-200 hover:underline">
              {t("new_password.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}