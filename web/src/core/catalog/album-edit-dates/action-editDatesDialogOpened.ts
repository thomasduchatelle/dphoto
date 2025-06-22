import {albumIdEquals, CatalogViewerState} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";
import {createAction} from "@light-state";

export const editDatesDialogOpened = createAction<CatalogViewerState>(
    "EditDatesDialogOpened",
    (current: CatalogViewerState) => {
        const {albumId: displayedAlbumId} = displayedAlbumSelector(current);

        const selectedAlbum = current.albums.find(album => displayedAlbumId && albumIdEquals(displayedAlbumId, album.albumId));

        if (!selectedAlbum) {
            return current;
        }

        const displayEndDate = new Date(selectedAlbum.end);
        if (displayEndDate.getHours() === 0 && displayEndDate.getMinutes() === 0 && displayEndDate.getSeconds() === 0) {
            displayEndDate.setDate(displayEndDate.getDate() - 1);
        }

        return {
            ...current,
            editDatesDialog: {
                albumId: selectedAlbum.albumId,
                albumName: selectedAlbum.name,
                startDate: selectedAlbum.start,
                endDate: displayEndDate,
                isLoading: false,
            },
        };
    }
);

export type EditDatesDialogOpened = ReturnType<typeof editDatesDialogOpened>;
