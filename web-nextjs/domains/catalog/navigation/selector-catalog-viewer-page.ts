import {Album, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";
import {albumRailWindow} from "./album-rail-window";
import {photoUrl} from "./photo-url";

export interface CatalogViewerPageSelection {
    albumsLoaded: boolean;
    albums: Album[];
    railAlbums: Album[];
    displayedAlbum: Album | undefined;
    previousAlbum: Album | undefined;
    nextAlbum: Album | undefined;
    medias: MediaWithinADay[];
    firstMediaHref: string | undefined;
    mediasLoaded: boolean;
    albumNotFound: boolean;
    error?: Error;
}

export function catalogViewerPageSelector(state: CatalogViewerState): CatalogViewerPageSelection {
    const displayedAlbumId = getDisplayedAlbumId(state);
    const displayedIndex = state.allAlbums.findIndex(album => albumIdEquals(album.albumId, displayedAlbumId));
    const displayedAlbum = displayedIndex >= 0 ? state.allAlbums[displayedIndex] : undefined;
    const displayedIndexInRail = state.albums.findIndex(album => albumIdEquals(album.albumId, displayedAlbumId));
    const firstMedia = state.medias[0]?.medias[0];

    return {
        albumsLoaded: state.albumsLoaded,
        albums: state.albums,
        railAlbums: albumRailWindow(state.albums, displayedIndexInRail),
        displayedAlbum,
        previousAlbum: displayedIndex >= 0 && displayedIndex < state.allAlbums.length - 1 ? state.allAlbums[displayedIndex + 1] : undefined,
        nextAlbum: displayedIndex > 0 ? state.allAlbums[displayedIndex - 1] : undefined,
        medias: state.medias,
        firstMediaHref: displayedAlbum && firstMedia ? photoUrl(displayedAlbum.albumId, firstMedia.id, {autoplay: true}) : undefined,
        mediasLoaded: state.mediasLoaded,
        albumNotFound: state.albumNotFound,
        error: state.error,
    };
}
