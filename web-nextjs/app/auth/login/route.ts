import {NextRequest, NextResponse} from 'next/server';
import {redirectUrl} from "@/libs/requests";
import {initiateAuthenticationFlow} from "@/libs/security";

export async function GET(request: NextRequest) {
    try {
        const redirectTo = await initiateAuthenticationFlow('/')
        return NextResponse.redirect(redirectTo.redirectTo)

    } catch (e) {
        console.error('Error during OAuth login initiation:', e);
        return NextResponse.redirect((await redirectUrl("/auth/error")));
    }
}
