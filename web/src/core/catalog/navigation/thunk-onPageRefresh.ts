import {Album, AlbumId, albumIdEquals, CatalogViewerState} from "../language";
import {AlbumsAndMediasLoaded, albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {MediaLoadFailed, mediaLoadFailed} from "./action-mediaLoadFailed";
import {MediasLoaded, mediasLoaded} from "./action-mediasLoaded";
import {NoAlbumAvailable} from "./action-noAlbumAvailable";
import {MediaPerDayLoader} from "./MediaPerDayLoader";
import {ThunkDeclaration} from "../../thunk-engine";
import {DPhotoApplication} from "../../application";
import {CatalogFactory} from "../catalog-factories";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

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
        private readonly dispatch: (action: MediaLoadFailed | AlbumsAndMediasLoaded | NoAlbumAvailable | MediasLoaded) => void,
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
                    this.dispatch(mediasLoaded({albumId, medias: medias.medias}))
                })
                .catch(error => {
                    this.dispatch(mediaLoadFailed({
                        selectedAlbum: allAlbums.find(a => albumIdEquals(a.albumId, albumId)),
                        error,
                    }))
                })
        }

        return Promise.resolve()
    }

    private loadSpecificAlbum = async (albumId: AlbumId): Promise<MediaLoadFailed | AlbumsAndMediasLoaded> => {
        const [albumsResp, mediasResp] = await Promise.allSettled([
            this.fetchAlbumsPort.fetchAlbums(),
            this.mediaPerDayLoader.loadMedias(albumId),
        ]);

        if (albumsResp.status === "rejected") {
            return Promise.reject(albumsResp.reason)

        } else if (mediasResp.status === "rejected") {
            const selectedAlbum = albumsResp.value.find((a: Album) => albumIdEquals(a.albumId, albumId))

            return mediaLoadFailed({
                albums: albumsResp.value,
                selectedAlbum: selectedAlbum,
                error: new Error(`failed to load medias of ${albumId}`, mediasResp.reason),
            })

        } else {
            const albums: Album[] = albumsResp.value
            const medias = mediasResp.value.medias

            const selectedAlbum = albums.find(a => albumIdEquals(a.albumId, albumId))
            return albumsAndMediasLoaded({
                albums: albums,
                medias,
                selectedAlbum,
            })
        }
    }

    private loadDefaultAlbum = async (): Promise<AlbumsAndMediasLoaded | NoAlbumAvailable | MediaLoadFailed> => {
        const albums = await this.fetchAlbumsPort.fetchAlbums();
        if (!albums || albums.length === 0) {
            return {type: 'noAlbumAvailable'} as NoAlbumAvailable
        }
        const selectedAlbum = albums[0];
        try {
            const medias = await this.mediaPerDayLoader.loadMedias(selectedAlbum.albumId);
            return albumsAndMediasLoaded({
                albums: albums,
                medias: medias.medias,
                selectedAlbum,
                redirectTo: selectedAlbum.albumId,
            });

        } catch (e: any) {
            return mediaLoadFailed({
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
