import "server-only"

import {NextRequest} from 'next/server';
import {authenticate,} from '@/libs/security';
import {buildRedirectResponse} from "@/libs/nextjs-cookies";


export async function GET(request: NextRequest) {
    const redirectTo = await authenticate(request);
    return buildRedirectResponse(redirectTo);
}
