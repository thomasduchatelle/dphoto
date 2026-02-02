import {NextRequest, NextResponse} from "next/server";
import {cookies} from "next/headers";

export type SetCookies = Record<string, CookieValue>;

export interface Redirection {
    redirectTo: URL;
    cookies?: SetCookies;
}

export interface ReadCookieStore {
    get: (name: string) => (string | undefined)
}

export interface CookieValue {
    value: string;
    maxAge: number; // use 0 to delete
}

export function appliesCookies(response: NextResponse, setCookies?: SetCookies): NextResponse {
    if(setCookies) {
        for (const [key, value] of Object.entries(setCookies)) {
            response.cookies.set(key, value.value, {
                maxAge: value.maxAge,
                path: '/',
                httpOnly: true,
                secure: true,
                sameSite: 'lax',
            });
        }
    }
    return response;
}


export function buildRedirectResponse(redirection: Redirection) {
    return appliesCookies(NextResponse.redirect(redirection.redirectTo), redirection.cookies);
}

export function newReadCookieStore(request: NextRequest): ReadCookieStore {
    return {
        get: (name: string) => {
            return request.cookies.get(name)?.value;
        }
    }
}

export async function newReadCookieStoreFromComponents(): Promise<ReadCookieStore> {
    const c = await cookies()
    return {
        get: (name: string) => {
            return c.get(name)?.value;
        }
    }
}