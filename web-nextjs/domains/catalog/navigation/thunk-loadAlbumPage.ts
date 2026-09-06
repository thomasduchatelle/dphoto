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

export class LoadAlbumPage {
    constructor(
        private readonly dispatch: (action: AlbumsAndMediasLoaded | MediaLoadFailed | NoAlbumAvailable) => void,
        private readonly port: LoadAlbumPagePort,
    ) {}

    loadAlbumPage = async (albumId: AlbumId): Promise<void> => {
        const [albumsResp, mediasResp] = await Promise.allSettled([
            this.port.fetchAlbums(),
            this.port.fetchMedias(albumId),
        ]);

        if (albumsResp.status === 'rejected') {
            return Promise.reject(albumsResp.reason);
        }

        const albums = albumsResp.value;

        if (!albums.some(a => albumIdEquals(a.albumId, albumId))) {
            this.dispatch(noAlbumAvailable(undefined));
            return;
        }

        if (mediasResp.status === 'rejected') {
            this.dispatch(mediaLoadFailed({
                albums,
                displayedAlbumId: albumId,
                error: new Error(`failed to load medias of ${JSON.stringify(albumId)}`, {cause: mediasResp.reason}),
            }));
            return;
        }

        this.dispatch(albumsAndMediasLoaded({
            albums,
            medias: mediasResp.value,
            mediasFromAlbumId: albumId,
        }));
    };
}

export const loadAlbumPageDeclaration: ThunkDeclaration<
    CatalogViewerState,
    Record<string, never>,
    (albumId: AlbumId) => Promise<void>,
    CatalogDispatch & { adapter: LoadAlbumPagePort }
> = {
    selector: () => ({}),
    factory: ({dispatch, adapter}) => new LoadAlbumPage(dispatch, adapter).loadAlbumPage,
};
