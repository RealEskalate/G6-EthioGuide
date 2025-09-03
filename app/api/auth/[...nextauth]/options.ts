import CredentialsProvider from "next-auth/providers/credentials";
import type { NextAuthOptions } from "next-auth";
import GoogleProvider from "next-auth/providers/google";


declare module "next-auth" {
  interface Session {
    accessToken?: string;
    refreshToken?: string;
    exp?: number;
    error?: string;
    errorDetails?: string;
    user?: {
      email?: string;
      role?: string;
    };
  }
}

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export const options: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID || "",
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || "",
    }),

    CredentialsProvider({
      name: "Credentials",
      credentials: {
        email: {
          label: "Email",
          type: "email",
          placeholder: "you@example.com",
        },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          throw new Error("Email and password are required");
        }

        try {
          const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
          });

          const result = await res.json();
          console.log("Authorize response:", result);

          if (res.ok && result.success && result.data) {
            const { access, refresh, role } = result.data;
            return {
              id: result.data.id || credentials.email,
              email: credentials.email,
              role: role || "applicant",
              accessToken: access,
              refreshToken: refresh,
            };
          }
          throw new Error(result.message || "Invalid email or password");
        } catch (err) {
          console.error("Login error:", err);
          throw new Error("Authentication failed");
        }
      },
    }),
  ],
  pages: {
    signIn: "/auth/login",
  },
  session: {
    strategy: "jwt",
    // maxAge: 24 * 60 * 60 * 60, // 1 day
  },
  callbacks: {
    async jwt({ token, user, trigger }) {
      if (user) {
        type UserWithRole = typeof user & {
          role?: string;
          accessToken?: string;
          refreshToken?: string;
        };
        const u = user as UserWithRole;
        const role = "role" in u ? u.role : undefined;
        console.log("Initial sign-in:", { email: u.email, role });
        token.user = {
          email: u.email,
          role,
        };
        token.email = u.email;
        token.role = role;
        token.accessToken = u.accessToken;
        token.refreshToken = u.refreshToken;
        token.exp = Math.floor(Date.now() / 1000) + 60000; // 15 minutes
      }

      if (
        trigger === "update" ||
        (token.exp && Date.now() > Number(token.exp) * 1000)
      ) {
        console.log("Token expired or update triggered, refreshing:", {
          trigger,
          exp: token.exp,
          currentTime: Math.floor(Date.now() / 1000),
          refreshToken: token.refreshToken,
        });
        try {
          if (!token.refreshToken) {
            throw new Error("No refresh token available");
          }
          const res = await fetch(`${API_URL}/auth/token/refresh`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token.refreshToken}`,
            },
            body: "",
          });
          const result = await res.json();
          console.log("Refresh response:", {
            status: res.status,
            success: result.success,
            data: result.data,
            message: result.message,
          });

          if (res.ok && result.success && result.data && result.data.access) {
            token.accessToken = result.data.access;
            // Reuse existing refresh token since response does not provide a new one
            token.exp = Math.floor(Date.now() / 1000) + 900; // 15 minutes
            delete token.errorDetails; // Clear any previous error details
            delete token.error; // Clear any previous error
          } else {
            throw new Error(result.message || "Token refresh failed");
          }
        } catch (err) {
          console.error("Token refresh error:", err);
          return {
            ...token,
            error: "RefreshAccessTokenError",
            errorDetails:
              typeof err === "object" && err !== null && "message" in err
                ? (err as { message: string }).message
                : String(err),
          };
        }
      }

      return token;
    },
    async session({ session, token }) {
      if (token.user) {
        session.user = token.user;
        session.accessToken = token.accessToken as string | undefined;
        session.refreshToken = token.refreshToken as string | undefined;
        session.exp = typeof token.exp === "number" ? token.exp : undefined;
        session.error =
          typeof token.error === "string" ? token.error : undefined;
        session.errorDetails =
          typeof token.errorDetails === "string"
            ? token.errorDetails
            : undefined;
      }
      console.log("Session updated:", {
        email: session.user?.email,
        exp: session.exp,
        error: session.error,
        errorDetails: session.errorDetails,
      });
      return session;
    },
  },
  secret: process.env.NEXTAUTH_SECRET,
  debug: true, // Enable for detailed logging
};
