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
import { registerSchema, RegisterFormData } from "@/lib/validation/register";
import { signIn } from "next-auth/react";
import { FaGoogle, FaEye, FaEyeSlash } from "react-icons/fa";
import { useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation"; 

export default function RegisterPage() {
  const { t } = useTranslation("auth");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const router = useRouter();

  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      fullName: "",
      username: "",
      email: "",
      phoneNumber: "",
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (data: RegisterFormData) => {
    try {
      const base = (process.env.NEXT_PUBLIC_API_URL || '').replace(/\/$/, '')
      const registerResponse = await fetch(`${base}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: data.fullName,
          username: data.username,
          email: data.email,
          password: data.password,
        }),
      });

      const registerResult: unknown = await registerResponse.json().catch(() => ({}));
      if (!registerResponse.ok) {
        const msg = (registerResult && typeof registerResult === 'object' && registerResult !== null && ('message' in registerResult || 'error' in registerResult))
          ? (String((registerResult as { message?: string; error?: string }).message || (registerResult as { error?: string }).error))
          : t('register.error');
        throw new Error(msg);
      }

      setSuccessMessage(t('register.success') || "Registration successful! Please check your email.");
      setTimeout(() => {
        setSuccessMessage("");
        router.push('/auth/check-email');
      }, 3000);
    } catch (err) {
      const message = (err instanceof Error && err.message) ? err.message : t('register.error');
      setErrorMessage(message);
    }
  };

  return (
    <div className="bg-neutral-light text-foreground flex flex-col items-center p-4 font-sans min-h-full">
      <div className="flex items-center gap-3">
        <Image
          src="/images/ethioguide-symbol.png"
          alt="EthioGuide Symbol"
          width={50}
          height={50}
          priority
        />
        <span className="text-gray-800 font-semibold text-3xl">EthioGuide</span>
      </div>
      <Card className="bg-background-light w-full max-w-md border-neutral mt-6 mb-8">
        <CardHeader>
          <div className="flex flex-col items-center space-y-4">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl sm:text-3xl font-bold text-center font-amharic bg-gradient-to-r from-[#2e4d57] to-[#1c3b2e] bg-clip-text text-transparent">
                {t("register.title")}
              </CardTitle>
            </div>
          </div>
          <p className="text-sm text-center text-[#2e4d57]/80 font-medium">
            {t("register.sub_title")}
          </p>
        </CardHeader>
        <CardContent className="px-6 pb-6">
          {successMessage && <p className="text-green-500 text-sm mb-4 text-center">{successMessage}</p>}
          {errorMessage && <p className="text-red-500 text-sm mb-4 text-center">{errorMessage}</p>}
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="fullName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("register.full_name")} *
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("register.full_name_placeholder")}
                        className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-red-500 text-xs" />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("register.username")} *
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("register.username_placeholder")}
                        className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-red-500 text-xs" />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("register.email")} *
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t("register.email_placeholder")}
                        className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-red-500 text-xs" />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="phoneNumber"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("register.phone_number")}
                    </FormLabel>
                    <FormControl>
                      <div className="flex">
                        <span className="inline-flex items-center px-3 text-sm text-foreground bg-white/80 backdrop-blur-sm border-2 border-r-0 border-[#a7b3b9]/50 rounded-l-xl h-12">
                          +251
                        </span>
                        <Input
                          placeholder={t("register.phone_number_placeholder")}
                          className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-l-none rounded-r-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90"
                          {...field}
                          value={field.value?.replace(/^\+251/, "") || ""}
                          onChange={(e) =>
                            field.onChange(
                              `+251${e.target.value.replace(/^\+251/, "")}`
                            )
                          }
                        />
                      </div>
                    </FormControl>
                    <FormMessage className="text-red-500 text-xs" />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("register.password")} *
                    </FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showPassword ? "text" : "password"}
                          placeholder={t("register.password_placeholder")}
                          className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90 pr-12"
                          {...field}
                        />
                        <button
                          type="button"
                          onClick={() => setShowPassword(!showPassword)}
                          className="absolute inset-y-0 right-0 flex items-center pr-4 text-[#2e4d57] hover:text-[#3a6a8d] transition-colors duration-200"
                          aria-label={
                            showPassword ? "Hide password" : "Show password"
                          }
                        >
                          {showPassword ? (
                            <FaEyeSlash className="h-5 w-5" />
                          ) : (
                            <FaEye className="h-5 w-5" />
                          )}
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
                      {t("register.confirm_password")} *
                    </FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showConfirmPassword ? "text" : "password"}
                          placeholder={t("register.confirm_password_placeholder")}
                          className="border-2 border-[#a7b3b9]/50 focus:border-[#3a6a8d] focus:ring-2 focus:ring-[#3a6a8d]/20 rounded-xl h-12 bg-white/80 backdrop-blur-sm transition-all duration-300 hover:bg-white/90 pr-12"
                          {...field}
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
                          {showConfirmPassword ? (
                            <FaEyeSlash className="h-5 w-5" />
                          ) : (
                            <FaEye className="h-5 w-5" />
                          )}
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
                {t("register.create_account")}
              </Button>
              <div className="flex items-center justify-center my-4">
                <div className="border-t border-[#a7b3b9] flex-grow"></div>
                <span className="px-4 text-[#2e4d57]/80 text-sm">or</span>
                <div className="border-t border-[#a7b3b9] flex-grow"></div>
              </div>
              <Button
                variant="outline"
                className="w-full border-neutral text-primary-dark hover:bg-secondary/20 rounded-md"
                onClick={() =>
                  signIn("google", { callbackUrl: "/api/auth/callback/google" })
                }
              >
                <FaGoogle className="h-4 w-4 mr-2" />
                {t("register.sign_in_with_google")}
              </Button>
            </form>
          </Form>
          <p className="mt-6 text-sm text-center text-[#2e4d57]/80">
            {t("register.have_account")}{" "}
            <Link href="/auth/login" className="text-[#3a6a8d] hover:text-[#2e4d57] font-semibold transition-colors duration-200 hover:underline">
              {t("register.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}