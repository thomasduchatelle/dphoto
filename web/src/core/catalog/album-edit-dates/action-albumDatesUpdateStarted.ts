import {CatalogViewerState} from "../language";

export interface AlbumDatesUpdateStarted {
    type: "AlbumDatesUpdateStarted";
}

export function albumDatesUpdateStarted(): AlbumDatesUpdateStarted {
    return {
        type: "AlbumDatesUpdateStarted",
    };
}

export function reduceAlbumDatesUpdateStarted(
    current: CatalogViewerState,
    action: AlbumDatesUpdateStarted,
): CatalogViewerState {
    if (!current.editDatesDialog) {
        return current;
    }

    return {
        ...current,
        editDatesDialog: {
            ...current.editDatesDialog,
            isLoading: true,
        },
    };
}

export function albumDatesUpdateStartedReducerRegistration(handlers: any) {
    handlers["AlbumDatesUpdateStarted"] = reduceAlbumDatesUpdateStarted as (
        state: CatalogViewerState,
        action: AlbumDatesUpdateStarted
    ) => CatalogViewerState;
}
