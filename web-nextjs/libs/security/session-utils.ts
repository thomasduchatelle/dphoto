import {cookies} from 'next/headers';
import {ACCESS_TOKEN_COOKIE, BackendSession, REFRESH_TOKEN_COOKIE} from './constants';
import {decodeJWTPayload, isOwnerFromJWT} from './jwt-utils';

const USER_INFO_COOKIE = 'dphoto-user-info';

interface UserInfo {
    name: string;
    email: string;
    picture?: string;
}

interface AccessTokenPayload {
    exp?: number;
    [key: string]: unknown;
}

export async function getBackendSession(): Promise<BackendSession | null> {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get(ACCESS_TOKEN_COOKIE)?.value;
    const userInfoCookie = cookieStore.get(USER_INFO_COOKIE)?.value;

    if (!accessToken) {
        return null;
    }

    const accessTokenPayload = decodeJWTPayload(accessToken) as AccessTokenPayload | null;
    const expiresAt = accessTokenPayload?.exp ? new Date(accessTokenPayload.exp * 1000) : new Date();

    let userInfo: UserInfo | null = null;
    if (userInfoCookie) {
        try {
            userInfo = JSON.parse(userInfoCookie);
        } catch (e) {
            console.error('Failed to parse user info cookie:', e);
        }
    }

    return {
        type: 'authenticated',
        accessToken: {
            accessToken: accessToken,
            expiresAt: expiresAt,
        },
        authenticatedUser: {
            name: userInfo?.name || '',
            email: userInfo?.email || '',
            picture: userInfo?.picture,
            isOwner: isOwnerFromJWT(accessToken),
        },
    };
}

export async function clearAuthCookies(): Promise<void> {
    const cookieStore = await cookies();
    cookieStore.set(ACCESS_TOKEN_COOKIE, '', { maxAge: 0, path: '/' });
    cookieStore.set(REFRESH_TOKEN_COOKIE, '', { maxAge: 0, path: '/' });
    cookieStore.set(USER_INFO_COOKIE, '', { maxAge: 0, path: '/' });
}
