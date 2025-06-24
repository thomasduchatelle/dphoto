import {albumIdEquals, CatalogViewerState} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";
import {createAction} from "src/libs/daction";
import {isRoundTime} from "../common/date-helper";

export const editDatesDialogOpened = createAction<CatalogViewerState>(
    "EditDatesDialogOpened",
    (current: CatalogViewerState) => {
        const {displayedAlbumId} = displayedAlbumSelector(current);

        const selectedAlbum = current.albums.find(album => displayedAlbumId && albumIdEquals(displayedAlbumId, album.albumId));

        if (!selectedAlbum) {
            return current;
        }

        const startDate = selectedAlbum.start;
        const endDate = new Date(selectedAlbum.end);

        const startAtDayStart = startDate.getUTCHours() === 0 && startDate.getUTCMinutes() === 0 && startDate.getUTCSeconds() === 0 && startDate.getUTCMilliseconds() === 0;
        const endAtDayEnd = endDate.getUTCHours() === 0 && endDate.getUTCMinutes() === 0 && endDate.getUTCSeconds() === 0 && endDate.getUTCMilliseconds() === 0;

        if (endAtDayEnd) {
            endDate.setDate(endDate.getDate() - 1);
        } else if (!isRoundTime(endDate)) { // Only subtract 1 minute if it's a precise time
            endDate.setUTCMinutes(endDate.getUTCMinutes() - 1);
        }

        return {
            ...current,
            editDatesDialog: {
                albumId: selectedAlbum.albumId,
                albumName: selectedAlbum.name,
                startDate: startDate,
                endDate: endDate,
                isLoading: false,
                startAtDayStart: startAtDayStart,
                endAtDayEnd: endAtDayEnd,
            },
        };
    }
);

export type EditDatesDialogOpened = ReturnType<typeof editDatesDialogOpened>;
