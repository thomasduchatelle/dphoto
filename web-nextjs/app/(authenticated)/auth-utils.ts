import {getAuthentication} from '@/libs/security';
import {newReadCookieStoreFromComponents} from '@/libs/nextjs-cookies';
import {CurrentUserInsight} from '@/domains/catalog/language';

export async function mustBeAuthenticated(): Promise<CurrentUserInsight> {
    const authentication = await getAuthentication(await newReadCookieStoreFromComponents());
    if (authentication.status !== 'authenticated') {
        throw new Error('Page requires authentication');
    }
    let {picture, isOwner}: { picture?: string; isOwner: boolean } = authentication.authenticatedUser;
    return {picture, isOwner};
}
