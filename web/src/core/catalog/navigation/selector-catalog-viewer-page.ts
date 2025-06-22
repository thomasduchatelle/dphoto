import {Album, AlbumFilterEntry, CatalogViewerState, MediaWithinADay} from "../language/catalog-state";
import {albumIdEquals} from "../language/utils-albumIdEquals";
import {AlbumId} from "../language";

export interface CatalogViewerPageSelection {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    albumsLoaded: boolean;
    albums: Album[];
    selectedAlbum: Album | undefined;
    medias: MediaWithinADay[];
    mediasLoaded: boolean;
    mediasLoadedFromAlbumId?: AlbumId;
    loadingMediasFor?: AlbumId;
    albumNotFound: boolean;
    error?: Error;
}

export function catalogViewerPageSelector(state: CatalogViewerState, selectedAlbumId?: AlbumId): CatalogViewerPageSelection {
    const selectedAlbum = state.albums.find(album => albumIdEquals(album.albumId, selectedAlbumId));

    return {
        albumFilter: state.albumFilter,
        albumFilterOptions: state.albumFilterOptions,
        albumsLoaded: state.albumsLoaded,
        albums: state.albums,
        selectedAlbum: selectedAlbum,
        medias: state.medias,
        mediasLoaded: state.mediasLoaded,
        mediasLoadedFromAlbumId: state.mediasLoadedFromAlbumId,
        loadingMediasFor: state.loadingMediasFor,
        albumNotFound: state.albumNotFound,
        error: state.error,
    };
}
