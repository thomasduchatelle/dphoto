import { CatalogViewerState, AlbumFilterCriterion } from "./catalog-state";
import { albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION } from "./catalog-common-modifiers";

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
    const filteredAlbums = current.allAlbums.filter(current.albumMatchCriterion ? current.albumMatchCriterion(action.criterion) : (album) => true);

    const allAlbumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? current.albumFilterOptions[0];

    return {
        ...current,
        albums: filteredAlbums,
        albumFilter: current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, action.criterion)) ?? allAlbumFilter,
    };
}
