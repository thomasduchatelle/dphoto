import "server-only"

import {NextRequest, NextResponse} from 'next/server';
import {authenticate,} from '@/libs/security';

export async function GET(request: NextRequest) {
    const redirectTo = await authenticate(request.nextUrl)
    return NextResponse.redirect(redirectTo.redirectTo)
}
