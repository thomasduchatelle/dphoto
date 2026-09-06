import {Album, AlbumId, albumIdEquals, CatalogViewerState, Media} from '../language';
import {albumsAndMediasLoaded, AlbumsAndMediasLoaded} from './action-albumsAndMediasLoaded';
import {mediaLoadFailed, MediaLoadFailed} from './action-mediaLoadFailed';
import {noAlbumAvailable, NoAlbumAvailable} from './action-noAlbumAvailable';
import {CatalogDispatch} from '../common/catalog-dispatch';
import {ThunkDeclaration} from '@/libs/dthunks';

export interface LoadAlbumPagePort {
    fetchAlbums(): Promise<Album[]>
    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

export async function loadAlbumPage(
    dispatch: (action: AlbumsAndMediasLoaded | MediaLoadFailed | NoAlbumAvailable) => void,
    port: LoadAlbumPagePort,
    albumId: AlbumId,
): Promise<void> {
    const [albumsResp, mediasResp] = await Promise.allSettled([
        port.fetchAlbums(),
        port.fetchMedias(albumId),
    ]);

    if (albumsResp.status === 'rejected') {
        return Promise.reject(albumsResp.reason);
    }

    const albums = albumsResp.value;

    if (!albums.some(a => albumIdEquals(a.albumId, albumId))) {
        dispatch(noAlbumAvailable(undefined));
        return;
    }

    if (mediasResp.status === 'rejected') {
        dispatch(mediaLoadFailed({
            albums,
            displayedAlbumId: albumId,
            error: new Error(`failed to load medias of ${JSON.stringify(albumId)}`, {cause: mediasResp.reason}),
        }));
        return;
    }

    dispatch(albumsAndMediasLoaded({
        albums,
        medias: mediasResp.value,
        mediasFromAlbumId: albumId,
    }));
}

export const loadAlbumPageDeclaration: ThunkDeclaration<
    CatalogViewerState,
    Record<string, never>,
    (albumId: AlbumId) => Promise<void>,
    CatalogDispatch & { adapter: LoadAlbumPagePort }
> = {
    selector: () => ({}),
    factory: ({dispatch, adapter}) => loadAlbumPage.bind(null, dispatch, adapter),
};
