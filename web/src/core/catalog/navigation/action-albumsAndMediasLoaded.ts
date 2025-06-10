import {Album, CatalogViewerState, MediaWithinADay, RedirectToAlbumIdAction} from "../language";

import {refreshFilters} from "../common/utils";

export interface AlbumsAndMediasLoaded extends RedirectToAlbumIdAction {
    type: 'albumsAndMediasLoaded'
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
}

export function albumsAndMediasLoaded(props: Omit<AlbumsAndMediasLoaded, "type">): AlbumsAndMediasLoaded {
    return {
        ...props,
        type: 'albumsAndMediasLoaded',
    };
}

export function reduceAlbumsAndMediasLoaded(
    current: CatalogViewerState,
    action: AlbumsAndMediasLoaded
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
    handlers["albumsAndMediasLoaded"] = reduceAlbumsAndMediasLoaded as (
        state: CatalogViewerState,
        action: AlbumsAndMediasLoaded
    ) => CatalogViewerState;
}
