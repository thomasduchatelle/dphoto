import {Album, AlbumFilterCriterion, AlbumId, MediaWithinADay} from "./catalog-state";

export type CatalogViewerAction =
    AlbumsAndMediasLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | MediasLoadedAction
    | AlbumsFilteredAction

export type AlbumsAndMediasLoadedAction = {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
}

export type MediaFailedToLoadAction = {
    type: 'MediaFailedToLoadAction'
    albums: Album[]
    selectedAlbum?: Album
    error: Error
}

export type NoAlbumAvailableAction = {
    type: 'NoAlbumAvailableAction'
}

export type StartLoadingMediasAction = {
    type: 'StartLoadingMediasAction'
    albumId: AlbumId
}

export type AlbumsFilteredAction = {
    type: 'AlbumsFilteredAction'
    criterion: AlbumFilterCriterion
    albumId?: AlbumId // albumId is set when the album is changing as well and the action behaves like a StartLoadingMediasAction
}

export function startLoadingMediasAction(albumId: AlbumId): StartLoadingMediasAction {
    return {type: 'StartLoadingMediasAction', albumId}
}

export type MediasLoadedAction = {
    type: 'MediasLoadedAction'
    albumId: AlbumId
    medias: MediaWithinADay[]
}

export function mediasLoadedAction(albumId: AlbumId, medias: MediaWithinADay[]): MediasLoadedAction {
    return {type: 'MediasLoadedAction', albumId, medias}
}

export function isCatalogViewerAction(arg: any): arg is CatalogViewerAction {
    return ['AlbumsAndMediasLoadedAction', 'MediaFailedToLoadAction', 'NoAlbumAvailableAction', 'StartLoadingMediasAction', 'MediasLoadedAction'].indexOf(arg.type) >= 0
}