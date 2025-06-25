import {albumIdEquals, CatalogViewerState, EditDatesDialog} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";
import {createAction} from "src/libs/daction";
import {convertFromModelToDisplayDate, isDateAtDayEnd, isDateAtDayStart} from "../date-range/date-helper";

export const editDatesDialogOpened = createAction<CatalogViewerState>(
    "EditDatesDialogOpened",
    (current: CatalogViewerState) => {
        const {displayedAlbumId} = displayedAlbumSelector(current);

        const selectedAlbum = current.albums.find(album => displayedAlbumId && albumIdEquals(displayedAlbumId, album.albumId));

        if (!selectedAlbum) {
            return current;
        }

        const startDate = selectedAlbum.start;
        const startAtDayStart = isDateAtDayStart(startDate);
        const endAtDayEnd = isDateAtDayEnd(selectedAlbum.end);
        const endDate = convertFromModelToDisplayDate(selectedAlbum.end, true, endAtDayEnd);

        const newDialog: EditDatesDialog = {
            type: "EditDatesDialog",
            albumId: selectedAlbum.albumId,
            albumName: selectedAlbum.name,
            startDate: startDate,
            endDate: endDate,
            isLoading: false,
            startAtDayStart: startAtDayStart,
            endAtDayEnd: endAtDayEnd,
        };

        return {
            ...current,
            dialog: newDialog,
        };
    }
);

export type EditDatesDialogOpened = ReturnType<typeof editDatesDialogOpened>;
