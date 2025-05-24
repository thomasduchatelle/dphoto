import {Album, albumMatchCriterion, CatalogViewerState} from "./catalog-state";
import {RedirectToAlbumIdAction} from "./catalog-actions";
import {albumIdEquals} from "./utils-albumIdEquals";
import {albumFilterAreCriterionEqual, DEFAULT_ALBUM_FILTER_ENTRY, generateAlbumFilterOptions} from "./catalog-common-modifiers";

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

/**
 * Reducer fragment for albumsLoadedAction.
 * Uses currentUser from the state.
 */
export function reduceAlbumsLoaded(
    current: CatalogViewerState,
    action: AlbumsLoadedAction,
): CatalogViewerState {
    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums)
    const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, current.albumFilter.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY

    let staging: CatalogViewerState = {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: action.albums.filter(albumMatchCriterion(current.albumFilter.criterion)),
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
