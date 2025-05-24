import { Album, CatalogViewerState } from "./catalog-state";
import { RedirectToAlbumIdAction } from "./catalog-actions";
import { generateAlbumFilterOptions } from "./catalog-reducer";
import { albumIdEquals } from "./utils-albumIdEquals";

export type AlbumsLoadedAction = RedirectToAlbumIdAction & {
    type: 'AlbumsLoadedAction'
    albums: Album[]
};

export function AlbumsLoadedAction(
    albums: Album[],
    redirectTo?: any // AlbumId | undefined
): AlbumsLoadedAction {
    return {
        type: 'AlbumsLoadedAction',
        albums,
        redirectTo,
    };
}

/**
 * Reducer fragment for AlbumsLoadedAction.
 * Uses currentUser from the state.
 */
export function reduceAlbumsLoaded(
    current: CatalogViewerState,
    action: Omit<AlbumsLoadedAction, "type">
): CatalogViewerState {
    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums);
    const albumFilter = albumFilterOptions.find(option =>
        option.criterion.selfOwned === current.albumFilter.criterion.selfOwned &&
        JSON.stringify(option.criterion.owners) === JSON.stringify(current.albumFilter.criterion.owners)
    ) ?? albumFilterOptions[0];

    let staging: CatalogViewerState = {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: action.albums.filter(
            album => {
                // Use the same logic as albumMatchCriterion, but avoid import for now
                if (albumFilter.criterion.selfOwned) {
                    // Owned by current user
                    return album.ownedBy === undefined;
                } else {
                    return albumFilter.criterion.owners.length === 0 ||
                        albumFilter.criterion.owners.includes(album.albumId.owner);
                }
            }
        ),
        error: undefined,
        albumsLoaded: true,
    };

    if (
        action.redirectTo &&
        !staging.albums.find(album => albumIdEquals(album.albumId, action.redirectTo))
    ) {
        const fallbackFilter = albumFilterOptions.find(option =>
            option.criterion.selfOwned === undefined &&
            option.criterion.owners.length === 0
        ) ?? albumFilterOptions[0];
        staging = {
            ...staging,
            albumFilter: fallbackFilter,
            albums: action.albums.filter(
                album => {
                    if (fallbackFilter.criterion.selfOwned) {
                        return album.ownedBy === undefined;
                    } else {
                        return fallbackFilter.criterion.owners.length === 0 ||
                            fallbackFilter.criterion.owners.includes(album.albumId.owner);
                    }
                }
            ),
        };
    }

    return staging;
}
