import { Album, MediaWithinADay, CatalogViewerState } from "./catalog-state";
import { RedirectToAlbumIdAction } from "./catalog-actions";
import { CurrentUserInsight, generateAlbumFilterOptions } from "./catalog-reducer";

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

// Reducer fragment for AlbumsAndMediasLoadedAction
export function reduceAlbumsAndMediasLoaded(
    current: CatalogViewerState,
    action: Omit<AlbumsAndMediasLoadedAction, "type">
): CatalogViewerState {
    // We need currentUser for generateAlbumFilterOptions, so we require the reducer to be partially applied with currentUser
    throw new Error("reduceAlbumsAndMediasLoaded requires currentUser: use makeReduceAlbumsAndMediasLoaded(currentUser)");
}

export function makeReduceAlbumsAndMediasLoaded(currentUser: CurrentUserInsight) {
    return (
        current: CatalogViewerState,
        action: Omit<AlbumsAndMediasLoadedAction, "type">
    ): CatalogViewerState => {
        const albumFilterOptions = generateAlbumFilterOptions(currentUser, action.albums);

        return {
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
    };
}
