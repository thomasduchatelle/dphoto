import {constructThunkFromDeclaration} from '@/libs/dthunks/server';
import {catalogThunks} from '@/domains/catalog/thunks';
import {initialCatalogState} from '@/domains/catalog/language/initial-catalog-state';
import {newServerAdapterFactory} from '@/domains/catalog/adapters/server-adapter-factory';
import {HomePageContent} from './_components/HomePageContent';
import {getCurrentAuthentication} from '@/libs/security';
import {newReadCookieStoreFromComponents} from '@/libs/nextjs-cookies';
import {CurrentUserInsight} from '@/domains/catalog/language/catalog-state';
import {ServerState} from "@/libs/dthunks/server/constructThunkFromDeclaration";

export default async function HomePage() {
    const authentication = await getCurrentAuthentication(await newReadCookieStoreFromComponents());

    const currentUser: CurrentUserInsight = {
        picture: authentication.status === 'authenticated' ? authentication.authenticatedUser.picture : undefined,
        isOwner: authentication.status === 'authenticated' ? authentication.authenticatedUser.isOwner : false,
    };

    const serverState = new ServerState(initialCatalogState(currentUser))

    const onPageRefresh = constructThunkFromDeclaration(
        catalogThunks.onPageRefresh,
        {
            adapter: newServerAdapterFactory(), dispatch: action => {
            }
        },
        serverState,
    );

    const catalogState = await onPageRefresh(undefined);

    return (
        <HomePageContent
            albums={catalogState.albums}
            isLoading={!catalogState.albumsLoaded}
            error={catalogState.error}
        />
    )
}
