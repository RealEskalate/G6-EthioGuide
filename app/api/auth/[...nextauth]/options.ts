import CredentialsProvider from "next-auth/providers/credentials";
import type { NextAuthOptions } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import { jwtDecode } from "jwt-decode";

declare module "next-auth" {
  interface Session {
    accessToken?: string;
    refreshToken?: string;
    exp?: number;
    error?: string;
    errorDetails?: string;
    user?: {
      id?: string;
      name?: string;
      email?: string;
      role?: string;
    };
  }
  interface User {
    id?: string;
    name?: string;
    email?: string;
    role?: string;
    accessToken?: string;
    refreshToken?: string;
  }
}

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export const options: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      checks: ["state"],
    }),

    CredentialsProvider({
      name: "Credentials",
      credentials: {
        provider: { label: "Provider", type: "text" },
        code: { label: "Auth Code", type: "text" },
        identifier: { label: "Email", type: "text" },
        password: { label: "Password", type: "password" },
        accessToken: { label: "Access Token", type: "text", optional: true },
        role: { label: "Role", type: "text", optional: true },
        id: { label: "ID", type: "text", optional: true },
        name: { label: "Name", type: "text", optional: true },
        email: { label: "Email", type: "text", optional: true },
      },
      async authorize(
        credentials: Record<
          "id" | "name" | "email" | "provider" | "code" | "identifier" | "password" | "accessToken" | "role",
          string
        > | undefined
      ) {
        // Case 1: Social login (already verified by backend)
        if (credentials?.accessToken && credentials?.role) {
          return {
            id: credentials.id ?? "social-user",
            name: credentials.name ?? "",
            email: credentials.email ?? "",
            role: credentials.role,
            accessToken: credentials.accessToken,
          };
        }

        // Case 2: Standard email/password login
        if (credentials?.identifier && credentials?.password) {
          if (!API_URL) {
            throw new Error("Missing NEXT_PUBLIC_API_URL for login");
          }
          const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              identifier: credentials.identifier,
              password: credentials.password,
            }),
          });

          const result = await res.json();
          if (res.ok && result.user && result.access_token) {
            return {
              id: result.user.id,
              name: result.user.name,
              email: result.user.email,
              role: result.user.role,
              accessToken: result.access_token,
            };
          }

          throw new Error(result.message || "Login failed");
        }

        // If neither case applies → invalid request
        throw new Error("Invalid credentials");
      },

      // async authorize(credentials) {
      //   if (!credentials?.identifier || !credentials?.password)
      //     throw new Error("Email/Username and password are required");

      //   const res = await fetch(`${API_URL}/auth/login`, {
      //     method: "POST",
      //     headers: { "Content-Type": "application/json" },
      //     body: JSON.stringify({
      //       identifier: credentials.identifier,
      //       password: credentials.password,
      //     }),
      //   });

      //   const result = await res.json();
      //   if (res.ok && result.user && result.access_token) {
      //     return {
      //       id: result.user.id,
      //       name: result.user.name,
      //       email: result.user.email,
      //       role: result.user.role,
      //       accessToken: result.access_token,
      //     };
      //   }

      //   throw new Error(result.message || "Login failed");
      // },
    }),
  ],
  pages: {
    signIn: "/auth/login",
  },
  session: {
    strategy: "jwt",
  },
  callbacks: {
    async redirect({ url, baseUrl }) {
      // Prevent open redirect vulnerabilities
      if (url.startsWith(baseUrl)) return url;
      if (url.startsWith("/")) return new URL(url, baseUrl).toString();
      return baseUrl;
    },
    async jwt({ token, user, trigger }) {
      if (user) {
        console.log("Initial sign-in:", {
          id: user.id,
          email: user.email,
          name: user.name,
          role: user.role,
        });
        token.user = {
          id: user.id,
          name: user.name,
          email: user.email,
          role: user.role,
        };
        token.email = user.email;
        token.role = user.role;
        token.accessToken = user.accessToken;
        // Parse accessToken expiry
        if (user.accessToken) {
          try {
            const decoded: { exp?: number } = jwtDecode(user.accessToken);
            token.exp = decoded.exp || Math.floor(Date.now() / 1000) + 900;
          } catch (err) {
            console.error("Failed to decode accessToken:", err);
            token.exp = Math.floor(Date.now() / 1000) + 900; // Fallback to 15 minutes
          }
        }
      }
      // Refresh access token if expired or on demand
      if (
        (trigger === "update" || (token.exp && Date.now() > Number(token.exp) * 1000)) &&
        API_URL && typeof token.accessToken === 'string' && token.accessToken
      ) {
        console.log("Token expired or update triggered, refreshing:", {
          trigger,
          exp: token.exp,
          currentTime: Math.floor(Date.now() / 1000),
        });
        try {
          const res = await fetch(`${API_URL}/auth/refresh`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token.accessToken}`,
            },
            body: "", // No body needed
          });
          const result = await res.json();
          console.log("Refresh response:", {
            status: res.status,
            access_token: result.access_token,
            message: result.message,
          });

          if (res.ok && result.access_token) {
            token.accessToken = result.access_token;
            // Parse new accessToken expiry
            try {
              const decoded: { exp?: number } = jwtDecode(result.access_token);
              token.exp = decoded.exp || Math.floor(Date.now() / 1000) + 900;
            } catch (err) {
              console.error("Failed to decode new accessToken:", err);
              token.exp = Math.floor(Date.now() / 1000) + 900; // Fallback
            }
            delete token.errorDetails;
            delete token.error;
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
      if (token.user && typeof token.user === "object") {
        const userObj = token.user as {
          id?: string;
          name?: string;
          email?: string;
          role?: string;
        };
        session.user = {
          id: userObj.id,
          name: userObj.name,
          email: userObj.email,
          role: userObj.role,
        };

        session.accessToken =
          typeof token.accessToken === "string" ? token.accessToken : undefined;
        session.exp = typeof token.exp === "number" ? token.exp : undefined;
        session.error =
          typeof token.error === "string" ? token.error : undefined;
        session.errorDetails =
          typeof token.errorDetails === "string"
            ? token.errorDetails
            : undefined;
      }
      console.log("Session updated:", {
        id: session.user?.id,
        email: session.user?.email,
        name: session.user?.name,
        role: session.user?.role,
        exp: session.exp,
        error: session.error,
        errorDetails: session.errorDetails,
      });
      return session;
    },
  },
  secret: process.env.NEXTAUTH_SECRET,
  debug: true,
};
