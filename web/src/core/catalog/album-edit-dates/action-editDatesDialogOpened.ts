import {AlbumId, CatalogViewerState} from "../language";

export interface EditDatesDialogOpenedAction {
    type: "EditDatesDialogOpened";
    albumId: AlbumId;
}

export function editDatesDialogOpenedAction(albumId: AlbumId): EditDatesDialogOpenedAction {
    return {
        albumId,
        type: "EditDatesDialogOpened",
    };
}

export function reduceEditDatesDialogOpened(
    current: CatalogViewerState,
    {albumId}: EditDatesDialogOpenedAction,
): CatalogViewerState {
    return {
        ...current,
        editDatesDialog: {
            albumId
        }
    };
}

export function editDatesDialogOpenedReducerRegistration(handlers: any) {
    handlers["EditDatesDialogOpened"] = reduceEditDatesDialogOpened as (
        state: CatalogViewerState,
        action: EditDatesDialogOpenedAction
    ) => CatalogViewerState;
}
