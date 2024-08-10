import {Album, AlbumId, MediaWithinADay} from "../catalog-react";

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

export type MediasLoadedAction = {
    type: 'MediasLoadedAction'
    albumId: AlbumId
    medias: MediaWithinADay[]
}

export type CatalogAction =
    AlbumsAndMediasLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | MediasLoadedAction
