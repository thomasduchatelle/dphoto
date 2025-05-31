import {CatalogViewerState} from "../catalog-state";

export interface StartDeleteAlbumAction {
    type: "StartDeleteAlbum";
}

export function startDeleteAlbumAction(): StartDeleteAlbumAction {
    return {
        type: "StartDeleteAlbum",
    };
}

export function reduceStartDeleteAlbum(
    current: CatalogViewerState,
    _: StartDeleteAlbumAction
): CatalogViewerState {
    if (!current.deleteDialog) {
        return current;
    }
    return {
        ...current,
        deleteDialog: {
            ...current.deleteDialog,
            isLoading: true,
            error: undefined,
        },
    };
}

export function startDeleteAlbumReducerRegistration(handlers: any) {
    handlers["StartDeleteAlbum"] = reduceStartDeleteAlbum as (
        state: CatalogViewerState,
        action: StartDeleteAlbumAction
    ) => CatalogViewerState;
}
