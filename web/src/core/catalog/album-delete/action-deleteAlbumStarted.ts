import {CatalogViewerState} from "../language";

export interface DeleteAlbumStarted {
    type: "deleteAlbumStarted";
}

export function deleteAlbumStarted(): DeleteAlbumStarted {
    return {
        type: "deleteAlbumStarted",
    };
}

export function reduceDeleteAlbumStarted(
    current: CatalogViewerState,
    _: DeleteAlbumStarted
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

export function deleteAlbumStartedReducerRegistration(handlers: any) {
    handlers["deleteAlbumStarted"] = reduceDeleteAlbumStarted as (
        state: CatalogViewerState,
        action: DeleteAlbumStarted
    ) => CatalogViewerState;
}
