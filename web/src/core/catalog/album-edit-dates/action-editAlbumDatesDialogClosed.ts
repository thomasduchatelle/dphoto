import {CatalogViewerState} from "../language";

export interface EditAlbumDatesDialogClosed {
    type: "EditAlbumDatesDialogClosed";
}

export function editAlbumDatesDialogClosed(): EditAlbumDatesDialogClosed {
    return {
        type: "EditAlbumDatesDialogClosed",
    };
}

export function reduceEditAlbumDatesDialogClosed(
    current: CatalogViewerState,
    action: EditAlbumDatesDialogClosed,
): CatalogViewerState {
    return {
        ...current,
        editAlbumDatesDialog: undefined,
    };
}

export function editAlbumDatesDialogClosedReducerRegistration(handlers: any) {
    handlers["EditAlbumDatesDialogClosed"] = reduceEditAlbumDatesDialogClosed as (
        state: CatalogViewerState,
        action: EditAlbumDatesDialogClosed
    ) => CatalogViewerState;
}
