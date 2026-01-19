import "server-only"

import {NextRequest} from 'next/server';
import {authenticate,} from '@/libs/security';
import {buildRedirectResponse, newReadCookieStore} from "@/libs/nextjs-cookies";


export async function GET(request: NextRequest) {
    const redirectTo = await authenticate(request.nextUrl, newReadCookieStore(request));
    return buildRedirectResponse(redirectTo);
}
