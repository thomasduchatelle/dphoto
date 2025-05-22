import {Album, AlbumFilterCriterion, AlbumId, MediaWithinADay, ShareError, Sharing} from "./catalog-state";

export type CatalogViewerAction =
    AlbumsAndMediasLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | MediasLoadedAction
    | AlbumsFilteredAction
    | AlbumsLoadedAction
    | OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction
    | CloseSharingModalAction
    | SharingModalErrorAction

// TODO Delete the type.
export type SharingModalAction =
    OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction
    | CloseSharingModalAction
    | SharingModalErrorAction

export type RedirectToAlbumIdAction = {
    redirectTo?: AlbumId
}

export type AlbumsAndMediasLoadedAction = RedirectToAlbumIdAction & {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
}

export type AlbumsLoadedAction = RedirectToAlbumIdAction & {
    type: 'AlbumsLoadedAction'
    albums: Album[]
}

export type MediaFailedToLoadAction = {
    type: 'MediaFailedToLoadAction'
    albums?: Album[]
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

export type AlbumsFilteredAction = RedirectToAlbumIdAction & {
    type: 'AlbumsFilteredAction'
    criterion: AlbumFilterCriterion
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

export function isRedirectToAlbumIdAction(arg: any): arg is RedirectToAlbumIdAction {
    return arg.redirectTo
}

export type OpenSharingModalAction = {
    type: "OpenSharingModalAction"
    albumId: AlbumId
}
export type AddSharingAction = {
    type: "AddSharingAction"
    sharing: Sharing
}
export type RemoveSharingAction = {
    type: "RemoveSharingAction"
    email: string
}
export type CloseSharingModalAction = {
    type: "CloseSharingModalAction"
}
export type SharingModalErrorAction = {
    type: "SharingModalErrorAction"
    error: ShareError
}
