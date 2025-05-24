import {AlbumFilterCriterion, albumMatchCriterion, CatalogViewerState} from "./catalog-state";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY} from "./catalog-common-modifiers";

export interface AlbumsFilteredAction {
    type: 'AlbumsFilteredAction'
    criterion: AlbumFilterCriterion
}

export function albumsFilteredAction(criterion: AlbumFilterCriterion): AlbumsFilteredAction {
    return { type: "AlbumsFilteredAction", criterion };
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
