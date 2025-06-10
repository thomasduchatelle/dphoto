import {CatalogViewerState} from "../language";

export interface AlbumDeleteFailed {
    type: "albumDeleteFailed";
    error: string;
}

export function albumDeleteFailed(error: string): AlbumDeleteFailed {
    if (!error || error.trim() === "") {
        throw new Error("AlbumDeleteFailed requires a non-empty error message");
    }
    return {
        error,
        type: "albumDeleteFailed",
    };
}

export function reduceAlbumDeleteFailed(
    current: CatalogViewerState,
    action: AlbumDeleteFailed,
): CatalogViewerState {
    if (!current.deleteDialog) {
        return current;
    }

    return {
        ...current,
        deleteDialog: {
            ...current.deleteDialog,
            error: action.error,
            isLoading: false,
        },
    };
}

export function albumDeleteFailedReducerRegistration(handlers: any) {
    handlers["albumDeleteFailed"] = reduceAlbumDeleteFailed as (
        state: CatalogViewerState,
        action: AlbumDeleteFailed
    ) => CatalogViewerState;
}
