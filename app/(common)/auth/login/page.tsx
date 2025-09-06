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
import { signIn } from "next-auth/react";
import { FaEye, FaEyeSlash, FaGoogle } from "react-icons/fa";
import { loginSchema, type LoginFormData } from "@/lib/validation/login";
import { useState } from "react";
import Image from "next/image";

export default function LoginPage() {
  const { t } = useTranslation("auth");
  const [showPassword, setShowPassword] = useState(false);

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      identifier: "",
      password: "",
    },
  });

  const onSubmit = async (data: LoginFormData) => {
    try {
      const result = await signIn("credentials", {
        redirect: false,
        identifier: data.identifier,
        password: data.password,
        callbackUrl: "/",
      });
      if (result?.error) {
        form.setError("root", { message: t("login.error") });
      } else if (result?.url) {
        window.location.href = result.url;
      }
    } catch {
      form.setError("root", { message: t("new_password.error") });
    }
  };

  return (
    <div className="bg-neutral-light text-foreground flex flex-col items-center p-5 font-sans min-h-full">
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
      <Card className="w-full max-w-md bg-white border border-gray-200 shadow-lg rounded-2xl overflow-hidden">
        <CardHeader className="pb-4">
          <div className="flex flex-col items-center space-y-3">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-2xl sm:text-3xl font-bold text-center font-amharic bg-gradient-to-r from-[#2e4d57] to-[#1c3b2e] bg-clip-text text-transparent">
                {t("login.title")}
              </CardTitle>
            </div>
          </div>
          <p className="text-sm text-center text-[#2e4d57]/80 font-medium">
            {t("register.sub_title")}
          </p>
        </CardHeader>

        <CardContent className="px-6 pb-6">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="identifier"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("login.identifier")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("login.identifier_placeholder")}
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
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic text-[#2e4d57] font-semibold">
                      {t("login.password")}
                    </FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showPassword ? "text" : "password"}
                          placeholder={t("login.password_placeholder")}
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

              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="remember-me"
                    className="h-4 w-4 rounded border-2 border-[#a7b3b9] text-[#3a6a8d] focus:ring-[#3a6a8d]/20 focus:ring-2"
                    style={{ accentColor: "#3a6a8d" }}
                  />
                  <label
                    htmlFor="remember-me"
                    className="text-sm text-[#2e4d57] font-medium"
                  >
                    {t("login.remember_me")}
                  </label>
                </div>
                <div className="text-sm">
                  <Link
                    href="/auth/reset-password"
                    className="text-[#3a6a8d] hover:text-[#2e4d57] font-medium transition-colors duration-200 hover:underline"
                  >
                    {t("login.change_password")}
                  </Link>
                </div>
              </div>

              {form.formState.errors.root && (
                <p className="text-red-500 text-sm bg-red-50 p-3 rounded-lg border border-red-200">
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
                {t("login.sign_in")}
              </Button>
            </form>
          </Form>
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
          <p className="mt-6 text-sm text-center text-[#2e4d57]/80">
          
            {t("login.new_to_ethioguide")}{" "}
            <Link
              href="/auth/register"
              className="text-[#3a6a8d] hover:text-[#2e4d57] font-semibold transition-colors duration-200 hover:underline"
            >
              {t("login.create_account")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
