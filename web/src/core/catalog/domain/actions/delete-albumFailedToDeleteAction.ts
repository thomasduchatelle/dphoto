import {CatalogViewerState} from "../catalog-state";

export interface AlbumFailedToDeleteAction {
    type: "AlbumFailedToDelete";
    error: string;
}

export function albumFailedToDeleteAction(props: Omit<AlbumFailedToDeleteAction, "type">): AlbumFailedToDeleteAction {
    if (!props.error || props.error.trim() === "") {
        throw new Error("AlbumFailedToDeleteAction requires a non-empty error message");
    }
    return {
        ...props,
        type: "AlbumFailedToDelete",
    };
}

export function reduceAlbumFailedToDelete(
    current: CatalogViewerState,
    action: AlbumFailedToDeleteAction,
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

export function albumFailedToDeleteReducerRegistration(handlers: any) {
    handlers["AlbumFailedToDelete"] = reduceAlbumFailedToDelete as (
        state: CatalogViewerState,
        action: AlbumFailedToDeleteAction
    ) => CatalogViewerState;
}
