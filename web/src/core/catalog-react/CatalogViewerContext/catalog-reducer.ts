import {albumIdEquals} from "../../catalog";
import {CatalogViewerAction, CatalogViewerState} from "./catalog-viewer-state";

export const initialCatalogState: CatalogViewerState = {
    albumNotFound: false,
    albums: [],
    medias: [],
    albumsLoaded: false,
    mediasLoaded: false,
}

export function catalogReducer(current: CatalogViewerState, action: CatalogViewerAction): CatalogViewerState {
    switch (action.type) {
        case "StartLoadingMediasAction":
            return {
                albumsLoaded: current.albumsLoaded,
                albums: current.albums,
                selectedAlbum: current.albums.find(album => albumIdEquals(album.albumId, action.albumId)),
                medias: [],
                error: undefined,
                loadingMediasFor: action.albumId,
                albumNotFound: false,
                mediasLoaded: false,
            }

        case "MediasLoadedAction":
            if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, action.albumId)) {
                // concurrency management - ignore if not the last album requested
                return current
            }

            const selectedAlbum = current.albums.find(album => albumIdEquals(album.albumId, action.albumId))
            return {
                ...current,
                loadingMediasFor: undefined,
                medias: action.medias,
                selectedAlbum,
                error: undefined,
                mediasLoaded: true,
            }

        case "NoAlbumAvailableAction":
            return {
                albumNotFound: true,
                albums: [],
                medias: [],
                albumsLoaded: true,
                mediasLoaded: true,
            }

        case "AlbumsAndMediasLoadedAction":
            return {
                albumNotFound: false,
                albums: action.albums,
                medias: action.media,
                selectedAlbum: action.selectedAlbum,
                albumsLoaded: true,
                mediasLoaded: true,
            }

        case "MediaFailedToLoadAction":
            return {
                albumNotFound: false,
                albums: action.albums ?? current.albums,
                medias: [],
                selectedAlbum: undefined,
                error: action.error,
                albumsLoaded: true,
                mediasLoaded: true,
            }

        default:
            return current
    }
}

