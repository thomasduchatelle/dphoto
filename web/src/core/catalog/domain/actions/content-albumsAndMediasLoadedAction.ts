import {Album, CatalogViewerState, MediaWithinADay, RedirectToAlbumIdAction} from "../catalog-state";

import {refreshFilters} from "../catalog-common-modifiers";

export interface AlbumsAndMediasLoadedAction extends RedirectToAlbumIdAction {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
}

export function albumsAndMediasLoadedAction(props: Omit<AlbumsAndMediasLoadedAction, "type">): AlbumsAndMediasLoadedAction {
    return {
        ...props,
        type: 'AlbumsAndMediasLoadedAction',
    };
}

export function reduceAlbumsAndMediasLoaded(
    current: CatalogViewerState,
    action: AlbumsAndMediasLoadedAction
): CatalogViewerState {
    const {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    return {
        currentUser: current.currentUser,
        albumNotFound: false,
        allAlbums: action.albums,
        albums: albums,
        mediasLoadedFromAlbumId: action.selectedAlbum?.albumId,
        medias: action.medias,
        albumsLoaded: true,
        mediasLoaded: true,
        albumFilterOptions,
        albumFilter: albumFilter,
    }
}

export function albumsAndMediasLoadedReducerRegistration(handlers: any) {
    handlers["AlbumsAndMediasLoadedAction"] = reduceAlbumsAndMediasLoaded as (
        state: CatalogViewerState,
        action: AlbumsAndMediasLoadedAction
    ) => CatalogViewerState;
}
