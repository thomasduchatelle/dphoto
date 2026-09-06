import {Metadata} from 'next';
import {notFound} from 'next/navigation';
import {serverSideThunk} from '@/libs/dthunks/server';
import {catalogThunks} from '@/domains/catalog/thunks';
import {initialCatalogState} from '@/domains/catalog/language/initial-catalog-state';
import {newServerSideRestCatalogAdapter} from '@/domains/catalog/adapters/server-adapter-factory';
import {AlbumPageContent} from './_components/AlbumPageContent';
import {mustBeAuthenticated} from '../../../auth-utils';

interface AlbumPageParams {
    ownerId: string;
    albumId: string;
}

export async function generateMetadata({params}: {
    params: Promise<AlbumPageParams>;
}): Promise<Metadata> {
    await params;
    return {
        title: 'DPhoto',
        description: 'Photo management application',
    };
}

export default async function AlbumPage({params}: { params: Promise<AlbumPageParams> }) {
    const {ownerId, albumId: folderName} = await params;

    const currentUser = await mustBeAuthenticated();
    const loadAlbumPage = serverSideThunk(
        catalogThunks.loadAlbumPage,
        {adapter: newServerSideRestCatalogAdapter()},
    );
    const catalogState = await loadAlbumPage(initialCatalogState(currentUser), {owner: ownerId, folderName});

    if (catalogState.error) {
        throw catalogState.error;
    }

    if (catalogState.albumNotFound) {
        notFound();
    }

    return <AlbumPageContent initialState={catalogState}/>;
}
