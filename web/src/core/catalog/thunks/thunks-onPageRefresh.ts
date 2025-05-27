import {
    Album,
    AlbumId,
    albumIdEquals,
    AlbumsAndMediasLoadedAction,
    catalogActions,
    CatalogViewerState,
    MediaFailedToLoadAction,
    MediaPerDayLoader,
    MediasLoadedAction,
    NoAlbumAvailableAction
} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";
import {DPhotoApplication} from "../../application";
import {CatalogFactory} from "../catalog-factories";
import {CatalogFactoryArgs} from "./catalog-factory-args";

export interface OnPageRefreshArgs {
    allAlbums: Album[]
    albumsLoaded: boolean
    mediasLoadedFromAlbumId?: AlbumId
    loadingMediasFor?: AlbumId
}

export interface FetchAlbumsPort {
    fetchAlbums(): Promise<Album[]>
}

export class OnPageRefresh {
    constructor(
        private readonly dispatch: (action: MediaFailedToLoadAction | AlbumsAndMediasLoadedAction | NoAlbumAvailableAction | MediasLoadedAction) => void,
        private readonly mediaPerDayLoader: MediaPerDayLoader,
        private readonly fetchAlbumsPort: FetchAlbumsPort
    ) {
    }

    onPageRefresh = async ({albumId, allAlbums, albumsLoaded, mediasLoadedFromAlbumId, loadingMediasFor}: OnPageRefreshArgs & {
        albumId?: AlbumId
    }): Promise<void> => {
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
        const [albumsResp, mediasResp] = await Promise.allSettled([
            this.fetchAlbumsPort.fetchAlbums(),
            this.mediaPerDayLoader.loadMedias(albumId),
        ]);

        if (albumsResp.status === "rejected") {
            return Promise.reject(albumsResp.reason)

        } else if (mediasResp.status === "rejected") {
            const selectedAlbum = albumsResp.value.find((a: Album) => albumIdEquals(a.albumId, albumId))

            return catalogActions.mediaFailedToLoadAction({
                albums: albumsResp.value,
                selectedAlbum: selectedAlbum,
                error: new Error(`failed to load medias of ${albumId}`, mediasResp.reason),
            })

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
        if (!albums || albums.length === 0) {
            return {type: 'NoAlbumAvailableAction'} as NoAlbumAvailableAction
        }
        const selectedAlbum = albums[0];
        try {
            const medias = await this.mediaPerDayLoader.loadMedias(selectedAlbum.albumId);
            return catalogActions.albumsAndMediasLoadedAction({
                albums: albums,
                medias: medias.medias,
                selectedAlbum,
                redirectTo: selectedAlbum.albumId,
            });

        } catch (e: any) {
            return catalogActions.mediaFailedToLoadAction({
                albums: albums,
                selectedAlbum,
                error: new Error(`failed to load medias of ${selectedAlbum.albumId}`, e),
            })
        }
    }
}

export const onPageRefreshDeclaration: ThunkDeclaration<
    CatalogViewerState,
    OnPageRefreshArgs,
    (albumId?: AlbumId) => Promise<void>,
    CatalogFactoryArgs
> = {
    factory: ({app, dispatch, partialState}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        const onPageRefreshInstance = new OnPageRefresh(
            dispatch,
            new MediaPerDayLoader(restAdapter),
            restAdapter
        );
        return (albumId?: AlbumId) => {
            const args = {
                ...partialState,
                albumId
            };
            return onPageRefreshInstance.onPageRefresh(args);
        };
    },
    selector: ({
                   allAlbums,
                   albumsLoaded,
                   mediasLoadedFromAlbumId,
                   loadingMediasFor
               }: CatalogViewerState): Omit<OnPageRefreshArgs, "albumId"> => ({
        allAlbums,
        albumsLoaded,
        mediasLoadedFromAlbumId,
        loadingMediasFor,
    })
}
