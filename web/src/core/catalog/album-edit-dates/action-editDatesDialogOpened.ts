import {albumIdEquals, CatalogViewerState} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";

export interface EditDatesDialogOpened {
    type: "EditDatesDialogOpened";
}

export function editDatesDialogOpened(): EditDatesDialogOpened {
    return {
        type: "EditDatesDialogOpened",
    };
}

export function reduceEditDatesDialogOpened(
    current: CatalogViewerState,
    _: EditDatesDialogOpened,
): CatalogViewerState {
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

export function editDatesDialogOpenedReducerRegistration(handlers: any) {
    handlers["EditDatesDialogOpened"] = reduceEditDatesDialogOpened as (
        state: CatalogViewerState,
        action: EditDatesDialogOpened
    ) => CatalogViewerState;
}
