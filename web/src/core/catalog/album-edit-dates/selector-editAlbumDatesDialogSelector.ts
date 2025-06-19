import {albumIdEquals, CatalogViewerState} from "../language";

export interface EditAlbumDatesDialogSelection {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
    isStartDateAtStartOfDay: boolean;
    isEndDateAtEndOfDay: boolean;
}

export const selectedEditAlbumDatesClosed: EditAlbumDatesDialogSelection = {
    isOpen: false,
    albumName: "",
    startDate: new Date(),
    endDate: new Date(),
    isStartDateAtStartOfDay: true,
    isEndDateAtEndOfDay: true,
};

export function editAlbumDatesDialogSelector(state: CatalogViewerState): EditAlbumDatesDialogSelection {
    if (!state.editAlbumDatesDialog) {
        return selectedEditAlbumDatesClosed;
    }

    const albumIdToEdit = state.editAlbumDatesDialog.albumId;
    const album = state.albums.find(a => albumIdEquals(a.albumId, albumIdToEdit));

    if (!album) {
        return selectedEditAlbumDatesClosed;
    }

    const inclusiveEndDate = new Date(album.end.getTime() - 1);

    return {
        isOpen: true,
        albumName: album.name,
        startDate: album.start,
        endDate: inclusiveEndDate,
        isStartDateAtStartOfDay: true,
        isEndDateAtEndOfDay: true,
    };
}
