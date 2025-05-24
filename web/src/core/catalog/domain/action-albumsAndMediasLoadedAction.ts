import {Album, CatalogViewerState, MediaWithinADay, RedirectToAlbumIdAction} from "./catalog-state";

import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY, generateAlbumFilterOptions} from "./catalog-common-modifiers";

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
    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums);

    return {
        currentUser: current.currentUser,
        albumNotFound: false,
        allAlbums: action.albums,
        albums: action.albums,
        mediasLoadedFromAlbumId: action.selectedAlbum?.albumId,
        medias: action.medias,
        albumsLoaded: true,
        mediasLoaded: true,
        albumFilterOptions,
        albumFilter: albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY
    }
}

export function albumsAndMediasLoadedReducerRegistration(handlers: any) {
    handlers["AlbumsAndMediasLoadedAction"] = reduceAlbumsAndMediasLoaded as (
        state: CatalogViewerState,
        action: AlbumsAndMediasLoadedAction
    ) => CatalogViewerState;
}
