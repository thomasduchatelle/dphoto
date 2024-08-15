import {Album, AlbumId, Media} from "./catalog-model";
import {AlbumsAndMediasLoadedAction, MediaFailedToLoadAction, NoAlbumAvailableAction} from "./catalog-actions";
import {MediaPerDayLoader} from "./SelectAlbumHandler";

export interface MediaViewLoaderQuery {
    albumId?: AlbumId
}

export interface FetchAlbumsPort {
    fetchAlbums(): Promise<Album[]>
}

export interface FetchAlbumMediasPort {
    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

// InitCatalogController is the controller used from loading page (or other external page)
export class CatalogViewerLoader {

    constructor(
        private readonly fetchAlbumsPort: FetchAlbumsPort,
        private readonly mediaPerDayLoader: MediaPerDayLoader,
    ) {
    }

    public loadInitialCatalog = (query: MediaViewLoaderQuery): Promise<MediaFailedToLoadAction | AlbumsAndMediasLoadedAction | NoAlbumAvailableAction> => {
        if (query.albumId) {
            return this.loadSpecificAlbum(query.albumId)
        }

        return this.loadDefaultAlbum()
    }

    private loadSpecificAlbum = (albumId: AlbumId): Promise<MediaFailedToLoadAction | AlbumsAndMediasLoadedAction> => {
        return Promise
            .allSettled([
                this.fetchAlbumsPort.fetchAlbums(),
                this.mediaPerDayLoader.loadMedias(albumId),
            ])
            .then(([albumsResp, mediasResp]) => {
                if (albumsResp.status === "rejected") {
                    return Promise.reject(albumsResp.reason)

                } else if (mediasResp.status === "rejected") {
                    return {
                        albums: albumsResp.value,
                        selectedAlbum: undefined,
                        error: new Error(`failed to load medias of ${albumId}`, mediasResp.reason),
                        type: 'MediaFailedToLoadAction',
                    } as MediaFailedToLoadAction

                } else {
                    const albums: Album[] = albumsResp.value
                    const medias = mediasResp.value.medias

                    const selectedAlbum = albums.find(a => a.albumId.owner === albumId.owner && a.albumId.folderName === albumId.folderName);
                    return {
                        albums: albums,
                        media: medias,
                        selectedAlbum,
                        type: 'AlbumsAndMediasLoadedAction',
                    } as AlbumsAndMediasLoadedAction
                }
            })
    }

    private loadDefaultAlbum = async (): Promise<AlbumsAndMediasLoadedAction | NoAlbumAvailableAction> => {
        const albums = await this.fetchAlbumsPort.fetchAlbums();
        if (!albums) {
            return {type: 'NoAlbumAvailableAction'} as NoAlbumAvailableAction
        }

        const selectedAlbum = albums[0];
        const medias = await this.mediaPerDayLoader.loadMedias(selectedAlbum.albumId);
        return ({
            type: 'AlbumsAndMediasLoadedAction',
            albums: albums,
            media: medias.medias,
            selectedAlbum,
        } as AlbumsAndMediasLoadedAction);
    }
}
