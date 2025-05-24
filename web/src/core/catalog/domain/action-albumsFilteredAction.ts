import {AlbumFilterCriterion, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "./catalog-state";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY} from "./catalog-common-modifiers";

export interface AlbumsFilteredAction extends RedirectToAlbumIdAction {
    type: 'AlbumsFilteredAction'
    criterion: AlbumFilterCriterion
}

export function albumsFilteredAction(props: Omit<AlbumsFilteredAction, "type">): AlbumsFilteredAction {
    return {...props, type: "AlbumsFilteredAction"};
}

export function reduceAlbumsFiltered(
    current: CatalogViewerState,
    action: AlbumsFilteredAction
): CatalogViewerState {
    const filteredAlbums = current.allAlbums.filter(albumMatchCriterion(action.criterion))

    const allAlbumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY

    return {
        ...current,
        albums: filteredAlbums,
        albumFilter: current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, action.criterion)) ?? allAlbumFilter,
    }
}

export function albumsFilteredReducerRegistration(handlers: any) {
    handlers["AlbumsFilteredAction"] = reduceAlbumsFiltered as (
        state: CatalogViewerState,
        action: AlbumsFilteredAction
    ) => CatalogViewerState;
}
