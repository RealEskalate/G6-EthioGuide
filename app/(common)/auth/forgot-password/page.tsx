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

  return (
    <div className="bg-neutral-light text-foreground flex flex-col items-center pt-4 pb-5 sm:pt-6 sm:pb-10 space-y-0 sm:space-y-1">
      <div className="flex items-center gap-2">
        <Image
          src="/images/ethioguide-symbol.png"
          alt="EthioGuide Symbol"
          width={40}
          height={40}
          className="h-10 w-10 sm:h-12 sm:w-12"
          priority
        />
        <span className="text-gray-800 font-semibold text-2xl sm:text-3xl">EthioGuide</span>
      </div>
      <Card className="bg-background-light w-full max-w-md border-neutral mb-0">
        <CardHeader className="px-4">
          <div className="flex flex-col items-center space-y-2">
            <div className="flex justify-center items-center w-full">
              <CardTitle className="text-xl sm:text-2xl font-bold text-center">
                {t("new_password.title")}
              </CardTitle>
            </div>
          </div>
          <p className="text-xs sm:text-sm text-center text-neutral-dark">
            {t("new_password.sub_title")}
          </p>
        </CardHeader>
        <CardContent className="pb-4">
          {!token ? (
            <p className="text-error text-xs sm:text-sm text-center">
              {t("new_password.invalid_token")}
            </p>
          ) : (
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-2 sm:space-y-3"
              >
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic text-sm sm:text-base">
                        {t("new_password.password")}
                      </FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showPassword ? "text" : "password"}
                            placeholder={t("new_password.password_placeholder")}
                            className="border-neutral focus:border-primary pr-10 text-sm sm:text-base"
                            {...field}
                            onBlur={() => form.trigger("password")}
                          />
                          <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute inset-y-0 right-0 flex items-center pr-3 text-neutral hover:text-primary"
                            aria-label={
                              showPassword ? "Hide password" : "Show password"
                            }
                          >
                            {showPassword ? (
                              <FaEyeSlash className="h-4 w-4 sm:h-5 sm:w-5" />
                            ) : (
                              <FaEye className="h-4 w-4 sm:h-5 sm:w-5" />
                            )}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-error text-xs sm:text-sm" />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-amharic text-sm sm:text-base">
                        {t("new_password.confirm_password")}
                      </FormLabel>
                      <FormControl>
                        <div className="relative">
                          <Input
                            type={showConfirmPassword ? "text" : "password"}
                            placeholder={t(
                              "new_password.confirm_password_placeholder"
                            )}
                            className="border-neutral focus:border-primary pr-10 text-sm sm:text-base"
                            {...field}
                            onBlur={() => form.trigger("confirmPassword")}
                          />
                          <button
                            type="button"
                            onClick={() =>
                              setShowConfirmPassword(!showConfirmPassword)
                            }
                            className="absolute inset-y-0 right-0 flex items-center pr-3 text-neutral hover:text-primary"
                            aria-label={
                              showConfirmPassword
                                ? "Hide confirm password"
                                : "Show confirm password"
                            }
                          >
                            {showConfirmPassword ? (
                              <FaEyeSlash className="h-4 w-4 sm:h-5 sm:w-5" />
                            ) : (
                              <FaEye className="h-4 w-4 sm:h-5 sm:w-5" />
                            )}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage className="text-error text-xs sm:text-sm" />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-primary to-primary-dark text-white hover:from-primary/90 hover:to-primary-dark/90 focus:ring-4 focus:ring-primary/50 rounded-md text-sm sm:text-base"
                >
                  {t("new_password.submit")}
                </Button>
                {form.formState.errors.root && (
                  <p className="text-error text-xs sm:text-sm text-center mt-2">
                    {form.formState.errors.root.message}
                  </p>
                )}
              </form>
            </Form>
          )}
          <p className="mt-2 text-xs sm:text-sm text-center text-neutral-dark">
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