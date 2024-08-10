import {Album, AlbumId, Media, MediaWithinADay} from "./catalog-model";
import {InternalError, UnrecoverableErrorAction} from "../application";
import {CatalogAction} from "./catalog-actions";

export interface MediaViewLoaderQuery {
    albumId?: AlbumId
}

export interface CatalogActionObserver {
    dispatch(action: CatalogAction): void
}

export interface FetchAlbumsPort {
    fetchAlbums(): Promise<Album[]>
}

export interface FetchAlbumMediasPort {
    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

export interface RedirectTo {
    albumId?: AlbumId
    albumPage: boolean
}

// InitCatalogController is the controller used from loading page (or other external page)
export class MediaViewLoader {

    constructor(
        private readonly fetchAlbumsPort: FetchAlbumsPort,
        private readonly fetchAlbumMediasPort: FetchAlbumMediasPort,
    ) {
    }

    public loadInitialCatalog = (query: MediaViewLoaderQuery): Promise<UnrecoverableErrorAction | CatalogAction> => {
        if (query.albumId) {
            return this.loadSpecificAlbum(query.albumId)
        }

        return this.loadDefaultAlbum()
    }

    private loadSpecificAlbum = (albumId: AlbumId): Promise<UnrecoverableErrorAction | CatalogAction> => {
        return Promise
            .allSettled([
                this.fetchAlbumsPort.fetchAlbums(),
                this.fetchAlbumMediasPort.fetchMedias(albumId),
            ])
            .then(([albumsResp, mediasResp]) => {
                if (albumsResp.status === "rejected") {
                    return {
                        type: 'unrecoverable-error',
                        error: new InternalError("failed to load albums", albumsResp.reason),
                    } as UnrecoverableErrorAction

                } else if (mediasResp.status === "rejected") {
                    return {
                        albums: albumsResp.value,
                        selectedAlbum: undefined,
                        error: new InternalError("failed to load medias", mediasResp.reason),
                        type: 'MediaFailedToLoadAction',
                    } as CatalogAction

                } else {
                    const albums = albumsResp.value
                    const medias = mediasResp.value

                    const selectedAlbum = albums.find(a => a.albumId.owner === albumId.owner && a.albumId.folderName === albumId.folderName);
                    return {
                        albums: albums,
                        media: groupByDay(medias),
                        selectedAlbum,
                        type: 'AlbumsAndMediasLoadedAction',
                    } as CatalogAction
                }
            })
            .catch(err => {
                // safeguard ... allSettled should never raise an error of this type
                return {
                    type: 'unrecoverable-error',
                    error: new InternalError("impossible error", err),
                } as UnrecoverableErrorAction
            })
    }

    private loadDefaultAlbum = (): Promise<UnrecoverableErrorAction | CatalogAction> => {
        return this.fetchAlbumsPort.fetchAlbums()
            .then(albums => {
                if (!albums) {
                    return {type: 'NoAlbumAvailableAction'} as CatalogAction
                }

                const selectedAlbum = albums[0];
                return this.fetchAlbumMediasPort.fetchMedias(selectedAlbum.albumId)
                    .then(medias => {
                        return {
                            type: 'AlbumsAndMediasLoadedAction',
                            albums: albums,
                            media: groupByDay(medias),
                            selectedAlbum,
                        } as CatalogAction
                    })
            })
            .catch(error => {
                return {
                    type: 'unrecoverable-error',
                    error: new InternalError("failed to load albums", error),
                } as UnrecoverableErrorAction
            });
    }
}


const groupByDay = (medias: Media[]): MediaWithinADay[] => {
    let result: MediaWithinADay[] = []

    medias.forEach(m => {
        const beginning = new Date(m.time)
        beginning.setHours(0, 0, 0, 0)

        if (result.length > 0 && result[0].day.getTime() === beginning.getTime()) {
            result[0].medias.push(m)
        } else {
            result = [{
                day: beginning,
                medias: [m],
            }, ...result]
        }
    })

    result.reverse()
    return result
}
