import {Album, CatalogViewerState, MediaWithinADay} from "./catalog-state";
import {RedirectToAlbumIdAction} from "./catalog-actions";
import {generateAlbumFilterOptions} from "./catalog-reducer";

export type AlbumsAndMediasLoadedAction = RedirectToAlbumIdAction & {
    type: 'AlbumsAndMediasLoadedAction'
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
};

export function AlbumsAndMediasLoadedAction(
    albums: Album[],
    medias: MediaWithinADay[],
    selectedAlbum?: Album,
    redirectTo?: any // AlbumId | undefined
): AlbumsAndMediasLoadedAction {
    return {
        type: 'AlbumsAndMediasLoadedAction',
        albums,
        medias,
        selectedAlbum,
        redirectTo,
    };
}

/**
 * Reducer fragment for AlbumsAndMediasLoadedAction.
 * Uses currentUser from the state.
 */
export function reduceAlbumsAndMediasLoaded(
    current: CatalogViewerState,
    action: Omit<AlbumsAndMediasLoadedAction, "type">
): CatalogViewerState {
    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums);

    return {
        ...current,
        albumNotFound: false,
        allAlbums: action.albums,
        albums: action.albums,
        mediasLoadedFromAlbumId: action.selectedAlbum?.albumId,
        medias: action.medias,
        albumsLoaded: true,
        mediasLoaded: true,
        albumFilterOptions,
        albumFilter: albumFilterOptions.find(option => option.criterion.selfOwned === undefined && option.criterion.owners.length === 0) ?? albumFilterOptions[0]
    };
}
