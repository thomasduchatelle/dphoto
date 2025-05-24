import {ActionObserver} from "./ActionObserver";
import {AlbumsAndMediasLoadedAction, catalogActions, MediaFailedToLoadAction, MediasLoadedAction, NoAlbumAvailableAction} from "./catalog-reducer-v2";
import {MediaPerDayLoader} from "./MediaPerDayLoader";
import {Album, AlbumId} from "./catalog-state";
import {albumIdEquals} from "./utils-albumIdEquals";

export interface PartialCatalogLoaderState {
    albumsLoaded: boolean,
    mediasLoadedFromAlbumId?: AlbumId
    loadingMediasFor?: AlbumId
    allAlbums: Album[]
}

export interface FetchAlbumsPort {
    fetchAlbums(): Promise<Album[]>
}

export class CatalogLoader {
    constructor(
        private readonly dispatch: ActionObserver<MediaFailedToLoadAction | AlbumsAndMediasLoadedAction | NoAlbumAvailableAction | MediasLoadedAction>,
        private readonly mediaPerDayLoader: MediaPerDayLoader,
        private readonly fetchAlbumsPort: FetchAlbumsPort,
        private readonly partialState: PartialCatalogLoaderState,
    ) {
    }

    onPageRefresh = (albumId?: AlbumId): Promise<void> => {
        const {mediasLoadedFromAlbumId, albumsLoaded, allAlbums, loadingMediasFor} = this.partialState

        if (!albumId) {
            if (!albumsLoaded) {
                return this.loadDefaultAlbum().then(this.dispatch)
            }

        } else if (!albumsLoaded && !albumIdEquals(mediasLoadedFromAlbumId, albumId) && !albumIdEquals(loadingMediasFor, albumId)) {
            return this.loadSpecificAlbum(albumId).then(this.dispatch)

        } else if (albumsLoaded && !albumIdEquals(mediasLoadedFromAlbumId, albumId) && !albumIdEquals(loadingMediasFor, albumId)) {
            return this.mediaPerDayLoader.loadMedias(albumId)
                .then(medias => {
                    this.dispatch(catalogActions.mediasLoadedAction({albumId, medias: medias.medias}))
                })
                .catch(error => {
                    this.dispatch(catalogActions.mediaFailedToLoadAction({
                        selectedAlbum: allAlbums.find(a => albumIdEquals(a.albumId, albumId)),
                        error,
                    }))
                })
        }

        return Promise.resolve()
    }

    private loadSpecificAlbum = async (albumId: AlbumId): Promise<MediaFailedToLoadAction | AlbumsAndMediasLoadedAction> => {
        const [albumsResp, mediasResp] = await Promise
            .allSettled([
                this.fetchAlbumsPort.fetchAlbums(),
                this.mediaPerDayLoader.loadMedias(albumId),
            ]);

        if (albumsResp.status === "rejected") {
            return Promise.reject(albumsResp.reason)

        } else if (mediasResp.status === "rejected") {
            const selectedAlbum = albumsResp.value.find(a => albumIdEquals(a.albumId, albumId))
            return {
                albums: albumsResp.value,
                selectedAlbum: selectedAlbum,
                error: new Error(`failed to load medias of ${albumId}`, mediasResp.reason),
                type: 'MediaFailedToLoadAction',
            } as MediaFailedToLoadAction

        } else {
            const albums: Album[] = albumsResp.value
            const medias = mediasResp.value.medias

            const selectedAlbum = albums.find(a => albumIdEquals(a.albumId, albumId))
            return catalogActions.albumsAndMediasLoadedAction({
                albums: albums,
                medias,
                selectedAlbum,
            })
        }
    }

    private loadDefaultAlbum = async (): Promise<AlbumsAndMediasLoadedAction | NoAlbumAvailableAction | MediaFailedToLoadAction> => {
        const albums = await this.fetchAlbumsPort.fetchAlbums();
        if (!albums) {
            return {type: 'NoAlbumAvailableAction'} as NoAlbumAvailableAction
        }

        const selectedAlbum = albums[0];
        try {
            const medias = await this.mediaPerDayLoader.loadMedias(selectedAlbum.albumId);
            return ({
                type: 'AlbumsAndMediasLoadedAction',
                albums: albums,
                medias: medias.medias,
                selectedAlbum,
                redirectTo: selectedAlbum.albumId,
            } as AlbumsAndMediasLoadedAction);

        } catch (e: any) {
            return {
                type: 'MediaFailedToLoadAction',
                albums: albums,
                selectedAlbum,
                error: new Error(`failed to load medias of ${selectedAlbum.albumId}`, e),
            } as MediaFailedToLoadAction
        }
    }
}
