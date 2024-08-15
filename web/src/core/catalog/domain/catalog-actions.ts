import {Album, AlbumId, MediaWithinADay} from "./catalog-model";

export type AlbumsAndMediasLoadedAction = {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    media: MediaWithinADay[]
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
