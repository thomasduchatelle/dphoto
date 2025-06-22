import {Album, AlbumId, albumIdEquals, CatalogViewerState, Media} from "../language";
import {albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {mediaLoadFailed} from "./action-mediaLoadFailed";
import {mediasLoaded} from "./action-mediasLoaded";
import {noAlbumAvailable} from "./action-noAlbumAvailable";
import {DPhotoApplication} from "../../application";
import {CatalogFactory} from "../catalog-factories";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ThunkDeclaration} from "src/libs/thunks";

export interface OnPageRefreshArgs {
    allAlbums: Album[]
    albumsLoaded: boolean
    mediasLoadedFromAlbumId?: AlbumId
    loadingMediasFor?: AlbumId
}

export interface FetchAlbumsAndMediasPort {
    fetchAlbums(): Promise<Album[]>

    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

export class OnPageRefresh {
    constructor(
        private readonly dispatch: (action: ReturnType<typeof mediaLoadFailed> | ReturnType<typeof albumsAndMediasLoaded> | ReturnType<typeof noAlbumAvailable> | ReturnType<typeof mediasLoaded>) => void,
        private readonly fetchAlbumsAndMediasPort: FetchAlbumsAndMediasPort
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
            return this.fetchAlbumsAndMediasPort.fetchMedias(albumId)
                .then(medias => {
                    this.dispatch(mediasLoaded({albumId, medias}))
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

    private loadSpecificAlbum = async (albumId: AlbumId): Promise<ReturnType<typeof mediaLoadFailed> | ReturnType<typeof albumsAndMediasLoaded>> => {
        const [albumsResp, mediasResp] = await Promise.allSettled([
            this.fetchAlbumsAndMediasPort.fetchAlbums(),
            this.fetchAlbumsAndMediasPort.fetchMedias(albumId),
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
            const medias = mediasResp.value

            const selectedAlbum = albums.find(a => albumIdEquals(a.albumId, albumId))
            return albumsAndMediasLoaded({
                albums: albums,
                medias,
                selectedAlbum,
            })
        }
    }

    private loadDefaultAlbum = async (): Promise<ReturnType<typeof albumsAndMediasLoaded> | ReturnType<typeof noAlbumAvailable> | ReturnType<typeof mediaLoadFailed>> => {
        const albums = await this.fetchAlbumsAndMediasPort.fetchAlbums();
        if (!albums || albums.length === 0) {
            return noAlbumAvailable()
        }
        const selectedAlbum = albums[0];
        try {
            const medias = await this.fetchAlbumsAndMediasPort.fetchMedias(selectedAlbum.albumId);
            return albumsAndMediasLoaded({
                albums: albums,
                medias: medias,
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
