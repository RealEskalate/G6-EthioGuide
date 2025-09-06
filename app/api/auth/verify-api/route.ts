// app/api/auth/verify/route.ts
import { NextRequest, NextResponse } from "next/server";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();

    if (!body.activatationToken) {
      return NextResponse.json(
        { message: "Activation token is missing" },
        { status: 400 }
      );
    }

    const backendResponse = await fetch(`${API_URL}/auth/verify`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const backendData = await backendResponse.json();
    console.log("Backend /auth/verify response:", backendData);
    return NextResponse.json(backendData, { status: backendResponse.status });
  } catch (error) {
    console.error("Error in /api/auth/verify:", error);
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    );
  }
}
