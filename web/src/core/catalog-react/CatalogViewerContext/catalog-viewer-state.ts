import {
    Album,
    AlbumId,
    AlbumsAndMediasLoadedAction,
    MediaFailedToLoadAction,
    MediasLoadedAction,
    MediaWithinADay,
    NoAlbumAvailableAction,
    StartLoadingMediasAction
} from "../../catalog";
import {Dispatch} from "react";

export type CatalogViewerAction =
    AlbumsAndMediasLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | MediasLoadedAction

export function isCatalogViewerAction(arg: any): arg is CatalogViewerAction {
    return ['AlbumsAndMediasLoadedAction', 'MediaFailedToLoadAction', 'NoAlbumAvailableAction', 'StartLoadingMediasAction', 'MediasLoadedAction'].indexOf(arg.type) >= 0
}

export interface CatalogViewerState {
    albumNotFound: boolean
    albums: Album[]
    selectedAlbum?: Album
    medias: MediaWithinADay[]
    error?: Error
    loadingMediasFor?: AlbumId
    albumsLoaded: boolean
    mediasLoaded: boolean
}

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    dispatch: Dispatch<CatalogViewerAction>
}