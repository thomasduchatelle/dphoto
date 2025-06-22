import {Album, AlbumFilterEntry, AlbumId, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";

export interface CatalogViewerPageSelection {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    albumsLoaded: boolean;
    albums: Album[];
    displayedAlbum: Album | undefined;
    medias: MediaWithinADay[];
    mediasLoaded: boolean;
    mediasLoadedFromAlbumId?: AlbumId;
    loadingMediasFor?: AlbumId;
    albumNotFound: boolean;
    error?: Error;
}

export function catalogViewerPageSelector(state: CatalogViewerState): CatalogViewerPageSelection {
    const {albumId: displayedAlbumId} = displayedAlbumSelector(state);
    const displayedAlbum = state.allAlbums.find(album => albumIdEquals(album.albumId, displayedAlbumId));

    return {
        albumFilter: state.albumFilter,
        albumFilterOptions: state.albumFilterOptions,
        albumsLoaded: state.albumsLoaded,
        albums: state.albums,
        displayedAlbum: displayedAlbum,
        medias: state.medias,
        mediasLoaded: state.mediasLoaded,
        mediasLoadedFromAlbumId: state.mediasLoadedFromAlbumId,
        loadingMediasFor: state.loadingMediasFor,
        albumNotFound: state.albumNotFound,
        error: state.error,
    };
}
