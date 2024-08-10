import {albumIdEquals, CatalogState} from "./catalog-model";
import {CatalogAction} from "../catalog";

export const initialCatalogState: CatalogState = {
    albumNotFound: false,
    albums: [],
    medias: [],
    fullyLoaded: false,
}

export function catalogReducer(current: CatalogState, action: CatalogAction): CatalogState {
    switch (action.type) {
        case "StartLoadingMediasAction":
            return {
                ...current,
                loadingMediasFor: action.albumId,
                albumNotFound: false,
                fullyLoaded: false,
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
                fullyLoaded: true,
            }

        case "NoAlbumAvailableAction":
            return {
                albumNotFound: true,
                albums: [],
                medias: [],
                fullyLoaded: true,
            }

        case "AlbumsAndMediasLoadedAction":
            return {
                albumNotFound: false,
                albums: action.albums,
                medias: action.media,
                selectedAlbum: action.selectedAlbum,
                fullyLoaded: true,
            }

        case "MediaFailedToLoadAction":
            return {
                albumNotFound: false,
                albums: action.albums ?? current.albums,
                medias: [],
                selectedAlbum: undefined,
                error: action.error,
                fullyLoaded: true,
            }

        default:
            return current
    }
}

