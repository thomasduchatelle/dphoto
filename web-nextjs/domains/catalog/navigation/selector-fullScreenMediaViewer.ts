import {AlbumId, CatalogViewerState, Media, MediaId} from "../language";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";

export interface FullScreenMediaViewerSelection {
    albumId?: AlbumId;
    medias: Media[];
    displayedIndex: number;
    mediaNotFound: boolean;
}

export function selectFullScreenMediaViewer(state: CatalogViewerState, mediaId: MediaId): FullScreenMediaViewerSelection {
    const albumId = getDisplayedAlbumId(state);
    const medias = state.medias.flatMap(group => group.medias);
    const displayedIndex = medias.findIndex(media => media.id === mediaId);

    return {
        albumId,
        medias,
        displayedIndex,
        mediaNotFound: displayedIndex < 0,
    };
}
