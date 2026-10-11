import {Album, AlbumId, albumIdEquals, CatalogViewerState, Media} from "../language";
import {AlbumsAndMediasLoaded, albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {AlbumsLoaded, albumsLoaded} from "./action-albumsLoaded";
import {MediaLoadFailed, mediaLoadFailed} from "./action-mediaLoadFailed";
import {mediasLoaded} from "./action-mediasLoaded";
import {NoAlbumAvailable, noAlbumAvailable} from "./action-noAlbumAvailable";
import {CatalogDispatch} from "../common/catalog-dispatch";
import {ThunkDeclaration} from "@/libs/dthunks";

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
        private readonly dispatch: (action: MediaLoadFailed | AlbumsAndMediasLoaded | AlbumsLoaded | NoAlbumAvailable | ReturnType<typeof mediasLoaded>) => void,
        private readonly fetchAlbumsAndMediasPort: FetchAlbumsAndMediasPort
    ) {
    }

    onPageRefresh = async ({albumId, albumsLoaded, mediasLoadedFromAlbumId, loadingMediasFor}: OnPageRefreshArgs & {
        albumId?: AlbumId
    }): Promise<void> => {
        if (!albumId) {
            if (!albumsLoaded) {
                return this.loadAlbumsOnly().then(this.dispatch)
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
                        displayedAlbumId: albumId,
                        error,
                    }))
                })
        }

        return Promise.resolve()
    }

    private loadSpecificAlbum = async (albumId: AlbumId): Promise<MediaLoadFailed | AlbumsAndMediasLoaded> => {
        const [albumsResp, mediasResp] = await Promise.allSettled([
            this.fetchAlbumsAndMediasPort.fetchAlbums(),
            this.fetchAlbumsAndMediasPort.fetchMedias(albumId),
        ]);

        if (albumsResp.status === "rejected") {
            return Promise.reject(albumsResp.reason)

        } else if (mediasResp.status === "rejected") {
            return mediaLoadFailed({
                albums: albumsResp.value,
                displayedAlbumId: albumId,
                error: new Error(`failed to load medias of ${albumId}`, mediasResp.reason),
            })

        } else {
            const albums: Album[] = albumsResp.value
            const medias = mediasResp.value

            return albumsAndMediasLoaded({
                albums: albums,
                medias,
                mediasFromAlbumId: albumId,
            })
        }
    }

    private loadAlbumsOnly = async (): Promise<AlbumsLoaded | NoAlbumAvailable> => {
        try {
            const albums = await this.fetchAlbumsAndMediasPort.fetchAlbums();

            if (!albums || albums.length === 0) {
                return noAlbumAvailable(undefined);
            }

            return albumsLoaded({albums});

        } catch (e) {
            console.log("onPageRefresh > loading albums failed", e);
            return noAlbumAvailable(e as Error);
        }
    }
}

export const onPageRefreshDeclaration: ThunkDeclaration<
    CatalogViewerState,
    OnPageRefreshArgs,
    (albumId?: AlbumId) => Promise<void>,
    CatalogDispatch & { adapter: FetchAlbumsAndMediasPort }
> = {
    factory: ({adapter, dispatch, partialState}) => {
        const onPageRefreshInstance = new OnPageRefresh(
            dispatch,
            adapter,
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
