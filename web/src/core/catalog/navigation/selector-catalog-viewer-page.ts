import {Album, AlbumFilterEntry, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";

export interface CatalogViewerPageSelection {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    albumsLoaded: boolean;
    albums: Album[];
    displayedAlbum: Album | undefined;
    medias: MediaWithinADay[];
    mediasLoaded: boolean;
    albumNotFound: boolean;
    error?: Error;
}

export function catalogViewerPageSelector(state: CatalogViewerState): CatalogViewerPageSelection {
    const displayedAlbumId = getDisplayedAlbumId(state);
    const displayedAlbum = state.allAlbums.find(album => albumIdEquals(album.albumId, displayedAlbumId));

    return {
        albumFilter: state.albumFilter,
        albumFilterOptions: state.albumFilterOptions,
        albumsLoaded: state.albumsLoaded,
        albums: state.albums,
        displayedAlbum: displayedAlbum,
        medias: state.medias,
        mediasLoaded: state.mediasLoaded,
        albumNotFound: state.albumNotFound,
        error: state.error,
    };
}
