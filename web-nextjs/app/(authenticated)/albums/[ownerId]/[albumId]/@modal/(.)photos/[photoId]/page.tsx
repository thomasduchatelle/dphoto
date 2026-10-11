import {serverSideThunk} from '@/libs/dthunks/server';
import {catalogThunks} from '@/domains/catalog/thunks';
import {initialCatalogState} from '@/domains/catalog/language/initial-catalog-state';
import {newServerSideRestCatalogAdapter} from '@/domains/catalog/adapters/server-adapter-factory';
import {mustBeAuthenticated} from '../../../../../../auth-utils';
import {PhotoViewerClient} from '../../../_components/PhotoViewerClient';

interface PhotoModalParams {
    ownerId: string;
    albumId: string;
    photoId: string;
}

export default async function PhotoModal({params}: { params: Promise<PhotoModalParams> }) {
    const {ownerId, albumId: folderName, photoId} = await params;
    const albumId = {owner: decodeURIComponent(ownerId), folderName: decodeURIComponent(folderName)};

    const currentUser = await mustBeAuthenticated();
    const loadAlbumPage = serverSideThunk(
        catalogThunks.loadAlbumPage,
        {adapter: newServerSideRestCatalogAdapter()},
    );
    const catalogState = await loadAlbumPage(initialCatalogState(currentUser), albumId);

    if (catalogState.error) {
        throw catalogState.error;
    }

    return (
        <PhotoViewerClient
            initialState={catalogState}
            mediaId={decodeURIComponent(photoId)}
            closeBehavior="back"
        />
    );
}
