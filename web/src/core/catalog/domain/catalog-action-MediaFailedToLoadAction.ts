import { Album, CatalogViewerState } from "./catalog-state";

/**
 * Action interface for MediaFailedToLoadAction
 */
export interface MediaFailedToLoadAction {
    type: 'MediaFailedToLoadAction'
    albums?: Album[]
    selectedAlbum?: Album
    error: Error
}

/**
 * Reducer fragment for MediaFailedToLoadAction.
 */
export function reduceMediaFailedToLoad(
    current: CatalogViewerState,
    action: Omit<MediaFailedToLoadAction, "type">
): CatalogViewerState {
    const allAlbums = action.albums ?? current.allAlbums;
    // generateAlbumFilterOptions and albumMatchCriterion are imported from catalog-reducer
    // but to avoid circular deps, import them in the index and pass as needed, or inline here if needed.
    // For now, we import them directly.
    // If you want to avoid import, you can pass them as params.
    // But for this migration, we import.
    // @ts-ignore
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    const { generateAlbumFilterOptions } = require("./catalog-reducer");
    // @ts-ignore
    const { albumMatchCriterion } = require("./catalog-state");

    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, allAlbums);
    const albumFilter = albumFilterOptions.find(
        (option: any) =>
            option.criterion.selfOwned === current.albumFilter.criterion.selfOwned &&
            JSON.stringify(option.criterion.owners) === JSON.stringify(current.albumFilter.criterion.owners)
    ) ?? albumFilterOptions[0];
    const albums = allAlbums.filter(albumMatchCriterion(albumFilter.criterion));

    return {
        currentUser: current.currentUser,
        allAlbums,
        albumFilterOptions,
        albumFilter,
        albums,
        albumNotFound: false,
        medias: [],
        error: action.error,
        albumsLoaded: true,
        mediasLoaded: true,
    };
}

/**
 * Action creator for MediaFailedToLoadAction
 */
export function mediaFailedToLoadAction(
    error: Error,
    albums?: Album[],
    selectedAlbum?: Album
): MediaFailedToLoadAction {
    return {
        type: "MediaFailedToLoadAction",
        error,
        albums,
        selectedAlbum,
    };
}
