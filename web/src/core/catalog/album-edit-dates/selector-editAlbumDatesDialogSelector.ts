import {AlbumId, CatalogViewerState} from "../../language";
import {albumIdEquals} from "../../language/utils-albumIdEquals";

export interface EditAlbumDatesDialogSelection {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
    isStartDateAtStartOfDay: boolean;
    isEndDateAtEndOfDay: boolean;
}

export function editAlbumDatesDialogSelector(state: CatalogViewerState): EditAlbumDatesDialogSelection {
    // Derived from `state.editAlbumDatesDialog` presence
    if (!state.editAlbumDatesDialog) {
        return {
            isOpen: false,
            albumName: "",
            startDate: new Date(),
            endDate: new Date(),
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        };
    }

    // From `state.editAlbumDatesDialog.albumId`
    const albumIdToEdit = state.editAlbumDatesDialog.albumId;
    // Derived from `state.albums` using `state.editAlbumDatesDialog.albumId`
    const album = state.albums.find(a => albumIdEquals(a.albumId, albumIdToEdit));

    // The "at the start/end of the day" checkboxes are always checked for this story.
    // The exclusive end date from the album needs to be adjusted to an inclusive end date for display.
    const inclusiveEndDate = album?.end ? new Date(album.end.getTime() - 1) : new Date();

    return {
        isOpen: true,
        albumName: album?.name || "", // From album.name
        startDate: album?.start || new Date(), // From album.start
        endDate: inclusiveEndDate, // From album.end (adjusted)
        isStartDateAtStartOfDay: true, // Hardcoded for this story
        isEndDateAtEndOfDay: true,     // Hardcoded for this story
    };
}
