import {CatalogViewerState} from "../language";

export interface DeleteAlbumDialogClosed {
    type: "deleteAlbumDialogClosed";
}

export function deleteAlbumDialogClosed(): DeleteAlbumDialogClosed {
    return {type: "deleteAlbumDialogClosed"};
}

export function reduceDeleteAlbumDialogClosed(
    state: CatalogViewerState,
    _action: DeleteAlbumDialogClosed
): CatalogViewerState {
    if (!state.deleteDialog) return state;
    return {
        ...state,
        deleteDialog: undefined,
    };
}

export function deleteAlbumDialogClosedReducerRegistration(handlers: any) {
    handlers["deleteAlbumDialogClosed"] = reduceDeleteAlbumDialogClosed as (
        state: CatalogViewerState,
        action: DeleteAlbumDialogClosed
    ) => CatalogViewerState;
}
