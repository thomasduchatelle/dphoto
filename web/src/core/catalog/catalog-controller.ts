// CatalogController is the controller used once on the page
import {AlbumId, Media, MediaWithinADay} from "./catalog-model";
import {CatalogAPIAdapter} from "../../apis/catalog";
import {Dispatch} from "react";
import {CatalogAction} from "./catalog-reducer";
import {InternalError, UnrecoverableErrorAction} from "../application";

export interface RedirectTo {
    albumId?: AlbumId
    albumPage: boolean
}

const noRedirectResponse: RedirectTo = {albumPage: false};
const redirectToAlbumsHomeResponse: RedirectTo = {albumPage: true};

// InitCatalogController is the controller used from loading page (or other external page)
export class InitialCatalogController {

    constructor(
        private readonly catalogAPIAdapter: CatalogAPIAdapter,
        private readonly dispatch: Dispatch<CatalogAction>,
        private readonly unrecoverableErrorDispatch: Dispatch<UnrecoverableErrorAction>,
    ) {
    }

    public loadInitialCatalog = (userEmail: string, owner?: string, folderName?: string): Promise<RedirectTo> => {
        if (owner && folderName) {
            return this.loadSpecificAlbum(userEmail, owner, folderName)
        }

        return this.loadDefaultAlbum(userEmail)
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

    private loadDefaultAlbum = (userEmail: string) => {
        return this.catalogAPIAdapter
            .fetchAlbums(userEmail)
            .then(albums => {
                if (!albums) {
                    this.dispatch({type: 'NoAlbumAvailableAction'})
                    return redirectToAlbumsHomeResponse
                }

                const selectedAlbum = albums[0];
                return this.catalogAPIAdapter.fetchMedias(selectedAlbum.albumId)
                    .then(medias => {
                        this.dispatch({
                            type: 'AlbumsAndMediasLoadedAction',
                            albums: albums,
                            media: groupByDay(medias),
                            selectedAlbum,
                        })
                    })
                    .then(() => ({albumId: selectedAlbum.albumId, albumPage: false}))
                    .catch(() => {
                        // todo dispatch

                        return redirectToAlbumsHomeResponse
                    })
            })
            .catch(error => {
                this.unrecoverableErrorDispatch({
                    type: 'unrecoverable-error',
                    error: new InternalError("failed to load albums", error),
                })

                return noRedirectResponse
            });
    }
}

export class CatalogController {
    constructor(
        private readonly catalogAPIAdapter: CatalogAPIAdapter,
        private readonly dispatch: Dispatch<CatalogAction>,
    ) {
    }

    public selectAlbum = (albumId: AlbumId): Promise<void> => {
        this.dispatch({type: "StartLoadingMediasAction", albumId: albumId})
        return this.catalogAPIAdapter.fetchMedias(albumId)
            .then(medias => {
                this.dispatch({type: 'MediasLoadedAction', albumId: albumId, medias: groupByDay(medias)})
            })
            .catch(error => this.dispatch({type: 'MediaFailedToLoadAction', error: error, albums: []}))
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
