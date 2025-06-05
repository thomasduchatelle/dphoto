import {Album, AlbumId, CatalogViewerState, RedirectToAlbumIdAction} from "../catalog-state";
import {albumIdEquals} from "../utils-albumIdEquals";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "./catalog-common-modifiers";

export interface AlbumDeletedAction extends RedirectToAlbumIdAction {
    type: "AlbumDeleted";
    albums: Album[];
    redirectTo?: AlbumId;
}

export function albumDeletedAction(props: Omit<AlbumDeletedAction, "type">): AlbumDeletedAction {
    return {
        ...props,
        type: "AlbumDeleted",
    };
}

export function reduceAlbumDeleted(
    current: CatalogViewerState,
    action: AlbumDeletedAction,
): CatalogViewerState {
    let {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    if (
        action.redirectTo &&
        !albums.some(album => albumIdEquals(album.albumId, action.redirectTo))
    ) {
        albumFilter =
            albumFilterOptions.find(option =>
                albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)
            ) ?? DEFAULT_ALBUM_FILTER_ENTRY;
        albums = action.albums
    }

    return {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: albums,
        error: undefined,
        albumsLoaded: true,
        deleteDialog: undefined,
    };
}

export function albumDeletedReducerRegistration(handlers: any) {
    handlers["AlbumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeletedAction
    ) => CatalogViewerState;
}
