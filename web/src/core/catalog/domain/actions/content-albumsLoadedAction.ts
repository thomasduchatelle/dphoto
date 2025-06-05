import {Album, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "../catalog-state";
import {albumIdEquals} from "../utils-albumIdEquals";
import {albumFilterAreCriterionEqual, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "./catalog-common-modifiers";

export interface AlbumsLoadedAction extends RedirectToAlbumIdAction {
    type: 'AlbumsLoadedAction'
    albums: Album[]
}

export function albumsLoadedAction(props: Omit<AlbumsLoadedAction, "type"> | Album[]): AlbumsLoadedAction {
    if (!props || Array.isArray(props)) {
        return albumsLoadedAction({albums: props});
    }

    return {
        ...props,
        type: 'AlbumsLoadedAction',
    };
}

export function reduceAlbumsLoaded(
    current: CatalogViewerState,
    action: AlbumsLoadedAction,
): CatalogViewerState {
    const {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    let staging: CatalogViewerState = {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: albums,
        error: undefined,
        albumsLoaded: true,
    }

    if (!staging.albums.find(album => albumIdEquals(album.albumId, action.redirectTo))) {
        const albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, DEFAULT_ALBUM_FILTER_ENTRY.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
        staging = {
            ...staging,
            albumFilter,
            albums: action.albums.filter(albumMatchCriterion(albumFilter.criterion)),
        }
    }

    return staging
}

export function albumsLoadedReducerRegistration(handlers: any) {
    handlers["AlbumsLoadedAction"] = reduceAlbumsLoaded as (
        state: CatalogViewerState,
        action: AlbumsLoadedAction
    ) => CatalogViewerState;
}
