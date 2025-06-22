import {AlbumId, albumIdEquals, CatalogViewerState, Media} from "../language";
import {createAction} from "src/light-state-lib";
import {groupByDay} from "./group-by-day";

interface MediasLoadedPayload {
    albumId: AlbumId
    medias: Media[]
}

export const mediasLoaded = createAction<CatalogViewerState, MediasLoadedPayload>(
    "mediasLoaded",
    (current: CatalogViewerState, {albumId, medias}: MediasLoadedPayload): CatalogViewerState => {
        if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, albumId)) {
            // concurrency management - ignore if not the last album requested
            return current
        }

        return {
            ...current,
            loadingMediasFor: undefined,
            mediasLoadedFromAlbumId: albumId,
            medias: groupByDay(medias),
            error: undefined,
            mediasLoaded: true,
            albumNotFound: false,
        }
    }
);

export type MediasLoaded = ReturnType<typeof mediasLoaded>;
