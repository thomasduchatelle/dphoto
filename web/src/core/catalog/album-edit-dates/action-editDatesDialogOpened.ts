import {AlbumId, CatalogViewerState} from "../language";

export interface EditDatesDialogOpened {
    type: "EditDatesDialogOpened";
    albumId: AlbumId;
}

export function editDatesDialogOpened(albumId: AlbumId): EditDatesDialogOpened {
    return {
        albumId,
        type: "EditDatesDialogOpened",
    };
}

export function reduceEditDatesDialogOpened(
    current: CatalogViewerState,
    {albumId}: EditDatesDialogOpened,
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
        action: EditDatesDialogOpened
    ) => CatalogViewerState;
}
