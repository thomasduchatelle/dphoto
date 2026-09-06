import {serverSideThunk} from '@/libs/dthunks/server';
import {catalogThunks} from '@/domains/catalog/thunks';
import {initialCatalogState} from '@/domains/catalog/language/initial-catalog-state';
import {newServerSideRestCatalogAdapter} from '@/domains/catalog/adapters/server-adapter-factory';
import {HomeContent} from './_components/HomeContent';
import {mustBeAuthenticated} from './auth-utils';

export default async function HomePage() {
    const currentUser = await mustBeAuthenticated();

    const onPageRefresh = serverSideThunk(
        catalogThunks.onPageRefresh,
        {adapter: newServerSideRestCatalogAdapter()},
    );
    const catalogState = await onPageRefresh(initialCatalogState(currentUser), undefined);

    if (catalogState.error) {
        throw catalogState.error;
    }

    return <HomeContent initialState={catalogState}/>;
}
