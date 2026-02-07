import {serverSideThunk} from '@/libs/dthunks/server';
import {catalogThunks} from '@/domains/catalog/thunks';
import {initialCatalogState} from '@/domains/catalog/language/initial-catalog-state';
import {HomePageContent} from './_components/HomePageContent';
import {getCurrentAuthentication} from '@/libs/security';
import {newReadCookieStoreFromComponents} from '@/libs/nextjs-cookies';
import {CurrentUserInsight} from '@/domains/catalog/language/catalog-state';
import {newServerSideRestCatalogAdapter} from "@/domains/catalog/adapters/server-adapter-factory";

export default async function HomePage() {
    const authentication = await getCurrentAuthentication(await newReadCookieStoreFromComponents());

    const currentUser: CurrentUserInsight = {
        picture: authentication.status === 'authenticated' ? authentication.authenticatedUser.picture : undefined,
        isOwner: authentication.status === 'authenticated' ? authentication.authenticatedUser.isOwner : false,
    };

    const onPageRefresh = serverSideThunk(
        catalogThunks.onPageRefresh,
        {
            adapter: newServerSideRestCatalogAdapter(),
        }
    )
    const catalogState = await onPageRefresh(initialCatalogState(currentUser), undefined);

    return (
        <HomePageContent
            albums={catalogState.albums}
            isLoading={!catalogState.albumsLoaded}
            error={catalogState.error}
        />
    )
}
