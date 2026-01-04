import {NextResponse} from 'next/server';
import * as cookie from 'cookie';

interface RedirectAssertion {
    url: string;
    params: Record<string, string | boolean>;
}

interface CookieAssertion {
    value?: string;
    maxAge?: number;
    httpOnly?: boolean;
    secure?: boolean;
    sameSite?: 'strict' | 'lax' | 'none';
    path?: string;
}

export function redirectionOf(response: NextResponse): RedirectAssertion {
    const location = response.headers.get('Location');
    if (!location) {
        throw new Error('No Location header found in response');
    }

    const url = new URL(location);
    const params: Record<string, string | boolean> = {};

    url.searchParams.forEach((value, key) => {
        // For empty values, just mark as present with true
        params[key] = value === '' ? true : value;
    });

    return {
        url: `${url.origin}${url.pathname}`,
        params,
    };
}

export function setCookiesOf(response: NextResponse): Record<string, CookieAssertion> {
    const setCookieHeaders = response.headers.getSetCookie();
    const result: Record<string, CookieAssertion> = {};

    for (const cookieStr of setCookieHeaders) {
        const parsed = cookie.parse(cookieStr);
        const [name, value] = Object.entries(parsed)[0];

        const assertion: CookieAssertion = {value};

        // Parse cookie attributes from the raw string
        if (cookieStr.includes('Max-Age=')) {
            const match = cookieStr.match(/Max-Age=(\d+)/i);
            if (match) {
                assertion.maxAge = parseInt(match[1], 10);
            }
        }

        if (cookieStr.includes('HttpOnly')) {
            assertion.httpOnly = true;
        }

        if (cookieStr.includes('Secure')) {
            assertion.secure = true;
        }

        if (cookieStr.includes('SameSite=')) {
            const match = cookieStr.match(/SameSite=(Strict|Lax|None)/i);
            if (match) {
                assertion.sameSite = match[1].toLowerCase() as 'strict' | 'lax' | 'none';
            }
        }

        if (cookieStr.includes('Path=')) {
            const match = cookieStr.match(/Path=([^;]+)/);
            if (match) {
                assertion.path = match[1].trim();
            }
        }

        result[name] = assertion;
    }

    return result;
}

