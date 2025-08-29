"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
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
import { LanguageSwitcher } from "@/components/language-switcher";
import { useState } from "react";
import Image from "next/image";


export default function RegisterPage() {
  const { t, i18n } = useTranslation("auth");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

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
    console.log("Form submitted:", data);
    // Example: await fetch('/api/register', { method: 'POST', body: JSON.stringify(data) });
  };

  // Debug translation loading
  // console.log('Current language:', i18n.language);
  // console.log('Auth translations:', i18n.getResourceBundle(i18n.language, 'auth'));

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
                {t("register.title")}
              </CardTitle>
              {/* <LanguageSwitcher /> */}
            </div>
          </div>
          <p className="text-sm text-neutral-dark ">{t("register.sub_title")}</p>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="fullName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">
                      {t("register.full_name")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("register.full_name_placeholder")}
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
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">
                      {t("register.username")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("register.username_placeholder")}
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
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">
                      {t("register.email")}
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t("register.email_placeholder")}
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
                name="phoneNumber"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="font-amharic">
                      {t("register.phone_number")}
                    </FormLabel>
                    <FormControl>
                      <div className="flex">
                        <span className="inline-flex items-center px-3 text-sm text-foreground bg-neutral/20 border border-r-0 border-neutral rounded-l-md">
                          +251
                        </span>
                        <Input
                          placeholder={t("register.phone_number_placeholder")}
                          className="border-neutral focus:border-primary rounded-l-none"
                          {...field}
                          value={field.value?.replace(/^\+251/, "") || ""} // Remove +251 for display
                          onChange={(e) =>
                            field.onChange(
                              `+251${e.target.value.replace(/^\+251/, "")}`
                            )
                          } // Prepend +251 on change
                        />
                      </div>
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
                    <FormLabel className="font-amharic">
                      {t("register.password")}
                    </FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showPassword ? "text" : "password"}
                          placeholder={t("register.password_placeholder")}
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
                    <FormLabel className="font-amharic">
                      {t("register.confirm_password")}
                    </FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showConfirmPassword ? "text" : "password"}
                          placeholder={t("register.confirm_password_placeholder")}
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
              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md"
              >
                {t("register.create_account")}
              </Button>
              <div className="flex items-center justify-center my-4">
                <div className="border-t border-neutral flex-grow"></div>
                <span className="px-4 text-neutral-dark text-sm">or</span>
                <div className="border-t border-neutral flex-grow"></div>
              </div>
              <Button
                variant="outline"
                className="w-full border-neutral text-primary-dark hover:bg-secondary/20 rounded-md"
                onClick={() => signIn("google", { callbackUrl: "/" })}
              >
                <FaGoogle className="h-4 w-4 mr-2" />
                {t("register.sign_in_with_google")}
              </Button>
            </form>
          </Form>
          <p className="mt-4 text-sm text-center text-neutral-dark">
            {t("register.have_account")}{" "}
            <Link href="/auth/login" className="text-primary hover:underline">
              {t("register.login_link")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}