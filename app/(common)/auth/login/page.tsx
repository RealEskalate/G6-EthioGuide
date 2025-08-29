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
import { FaEye, FaEyeSlash } from "react-icons/fa";
import { LanguageSwitcher } from "@/components/language-switcher";
import { loginSchema, LoginFormData } from "@/lib/validation/login";
import { useState } from "react";
import Image from "next/image";

export default function LoginPage() {
  const { t, i18n } = useTranslation("auth");
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
    } catch (error) {
      form.setError("root", { message: t("login.error") });
    }
  };

  // Debug translation loading
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
            <div className="flex justify-start items-center w-full">
              <CardTitle className="text-2xl font-bold text-center font-amharic">
                {t("login.title")}
              </CardTitle>
              {/* <LanguageSwitcher /> */}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="identifier"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">{t("login.identifier")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("login.identifier_placeholder")}
                        className="border-neutral focus:border-primary"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className="text-error" />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">{t("login.password")}</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showPassword ? "text" : "password"}
                          placeholder={t("login.password_placeholder")}
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
              <div className="flex justify-between items-center">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="remember-me"
                    className="h-4 w-4 text-primary focus:ring-primary border-neutral rounded"
                  />
                  <label htmlFor="remember-me" className="text-sm text-black ">
                    {t("login.remember_me")}
                  </label>
                </div>
                <div className="text-sm text-neutral-dark">
                  <Link href="/auth/reset-password" className="text-primary hover:underline">
                    {t("login.change_password")}
                  </Link>
                </div>
              </div>
              {form.formState.errors.root && (
                <p className="text-error text-sm">{form.formState.errors.root.message}</p>
              )}
              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md"
              >
                {t("login.sign_in")}
              </Button>
            </form>
          </Form>
          <p className="mt-2 text-sm text-center text-neutral-dark">
            {t("login.new_to_ethioguide")}{" "}
            <Link href="/auth/register" className="text-primary hover:underline">
              {t("login.create_account")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}