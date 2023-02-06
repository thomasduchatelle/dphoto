import {Album, AlbumId, albumIdEquals, CatalogState, MediaWithinADay} from "./catalog-model";

type AlbumsAndMediasLoadedAction = {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    media: MediaWithinADay[]
    selectedAlbum?: Album
}

type MediaFailedToLoadAction = {
    type: 'MediaFailedToLoadAction'
    albums: Album[]
    selectedAlbum?: Album
    error: Error
}

type NoAlbumAvailableAction = {
    type: 'NoAlbumAvailableAction'
}

type StartLoadingMediasAction = {
    type: 'StartLoadingMediasAction'
    albumId: AlbumId
}

type MediasLoadedAction = {
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

export const initialCatalogState: CatalogState = {
    albumNotFound: false, albums: [], medias: [],
}

export function catalogReducer(current: CatalogState, action: CatalogAction): CatalogState {
    switch (action.type) {
        case "StartLoadingMediasAction":
            return {
                ...current,
                loadingMediasFor: action.albumId,
                albumNotFound: false,
            }

        case "MediasLoadedAction":
            if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, action.albumId)) {
                // concurrency management - ignore if not the last album requested
                return current
            }

            const selectedAlbum = current.albums.find(album => albumIdEquals(album.albumId, action.albumId))
            return {
                ...current,
                loadingMediasFor: undefined,
                medias: action.medias,
                selectedAlbum,
                error: undefined,
            }

        case "NoAlbumAvailableAction":
            return {
                albumNotFound: true,
                albums: [],
                medias: [],
            }

        case "AlbumsAndMediasLoadedAction":
            return {
                albumNotFound: false,
                albums: action.albums,
                medias: action.media,
                selectedAlbum: action.selectedAlbum,
            }

        case "MediaFailedToLoadAction":
            return {
                albumNotFound: false,
                albums: action.albums ?? current.albums,
                medias: [],
                selectedAlbum: undefined,
                error: action.error,
            }

    }
    return current
}

