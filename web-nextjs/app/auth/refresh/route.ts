import "server-only"

import {NextRequest, NextResponse} from 'next/server';
import {refreshSession} from "@/libs/security";
import {appliesCookies, newReadCookieStore} from "@/libs/nextjs-cookies";

export async function POST(request: NextRequest) {
    const result = await refreshSession(newReadCookieStore(request));
    const status = result.success ? 204 : 401;
    return appliesCookies(new NextResponse(null, {status}), result.cookies);
}
