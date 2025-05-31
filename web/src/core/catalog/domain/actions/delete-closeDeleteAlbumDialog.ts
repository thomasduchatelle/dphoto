import {CatalogViewerState} from "../catalog-state";

export interface CloseDeleteAlbumDialogAction {
    type: "CloseDeleteAlbumDialog";
}

export function closeDeleteAlbumDialogAction(): CloseDeleteAlbumDialogAction {
    return { type: "CloseDeleteAlbumDialog" };
}

export function reduceCloseDeleteAlbumDialog(
    state: CatalogViewerState,
    _action: CloseDeleteAlbumDialogAction
): CatalogViewerState {
    if (!state.deleteDialog) return state;
    return {
        ...state,
        deleteDialog: undefined,
    };
}

export function closeDeleteAlbumDialogReducerRegistration(handlers: any) {
    handlers["CloseDeleteAlbumDialog"] = reduceCloseDeleteAlbumDialog as (
        state: CatalogViewerState,
        action: CloseDeleteAlbumDialogAction
    ) => CatalogViewerState;
}
