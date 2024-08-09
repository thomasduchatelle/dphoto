import {InternalError} from "../application";
import {Album, AlbumId, Media, RedirectTo} from "../catalog-react";

export interface FetchAlbumsPort {
    fetchAlbums(email: string): Promise<Album[]>
}

export interface FetchMediasPort {
    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

export interface LoaderMediaViewObserver {

}

export class LoaderMediaView {
    constructor(
        private readonly fetchAlbumsPort: FetchAlbumsPort,
        private readonly fetchMediasPort: FetchMediasPort,
        private readonly observers: LoaderMediaViewObserver[],
    ) {
    }


    private loadSpecificAlbum = (userEmail: string, owner: string, folderName: string): Promise<RedirectTo> => {
        return Promise
            .allSettled([
                this.catalogAPIAdapter.fetchAlbums(userEmail),
                this.catalogAPIAdapter.fetchMedias({owner, folderName}),
            ])
            .then(([albumsResp, mediasResp]) => {
                if (albumsResp.status === "rejected") {
                    this.unrecoverableErrorDispatch({
                        type: 'unrecoverable-error',
                        error: new InternalError("failed to load albums", albumsResp.reason),
                    })

                } else if (mediasResp.status === "rejected") {
                    this.dispatch({
                        albums: albumsResp.value,
                        selectedAlbum: undefined,
                        error: new InternalError("failed to load medias", mediasResp.reason),
                        type: 'MediaFailedToLoadAction',
                    })

                } else {
                    const albums = albumsResp.value
                    const medias = mediasResp.value

                    const selectedAlbum = albums.find(a => a.albumId.owner === owner && a.albumId.folderName === folderName);
                    this.dispatch({
                        albums: albums,
                        media: groupByDay(medias),
                        selectedAlbum,
                        type: 'AlbumsAndMediasLoadedAction',
                    })
                }
            })
            .then(() => noRedirectResponse)
            .catch(err => {
                // safe guard ... allSettled should never raise an error of this type
                this.unrecoverableErrorDispatch({
                    type: 'unrecoverable-error',
                    error: new InternalError("impossible error", err),
                })

                return noRedirectResponse
            })
    }
}