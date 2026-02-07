import {FetchCatalogAdapter} from "@/domains/catalog/adapters/api";
import {newReadCookieStoreFromComponents} from "@/libs/nextjs-cookies";
import {loadSession} from "@/libs/security/backend-store";
import {basePath, newOriginFromHeaders} from "@/libs/requests";

export function newServerSideRestCatalogAdapter(): FetchCatalogAdapter {
    return new FetchCatalogAdapter(
        async (): Promise<string | undefined> => {
            const cookieStore = await newReadCookieStoreFromComponents();
            const session = loadSession(cookieStore);
            return session.accessToken;
        },
        async () => {
            const url = await newOriginFromHeaders().getCurrentUrl();
            const prefix = process.env.NODE_ENV === 'development' ? basePath : "";
            return `${url.origin}${prefix}/api/v1`;
        }
    );
}
