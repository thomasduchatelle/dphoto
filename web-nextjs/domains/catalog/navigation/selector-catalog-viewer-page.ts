import {Album, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";

export interface CatalogViewerPageSelection {
    albumsLoaded: boolean;
    albums: Album[];
    displayedAlbum: Album | undefined;
    previousAlbum: Album | undefined;
    nextAlbum: Album | undefined;
    medias: MediaWithinADay[];
    mediasLoaded: boolean;
    albumNotFound: boolean;
    error?: Error;
}

export function catalogViewerPageSelector(state: CatalogViewerState): CatalogViewerPageSelection {
    const displayedAlbumId = getDisplayedAlbumId(state);
    const displayedIndex = state.allAlbums.findIndex(album => albumIdEquals(album.albumId, displayedAlbumId));
    const displayedAlbum = displayedIndex >= 0 ? state.allAlbums[displayedIndex] : undefined;

    return {
        albumsLoaded: state.albumsLoaded,
        albums: state.albums,
        displayedAlbum,
        previousAlbum: displayedIndex > 0 ? state.allAlbums[displayedIndex - 1] : undefined,
        nextAlbum: displayedIndex >= 0 && displayedIndex < state.allAlbums.length - 1 ? state.allAlbums[displayedIndex + 1] : undefined,
        medias: state.medias,
        mediasLoaded: state.mediasLoaded,
        albumNotFound: state.albumNotFound,
        error: state.error,
    };
}
