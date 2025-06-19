import {AlbumId, CatalogViewerState} from "../language";

export interface EditAlbumDatesDialogOpened {
    type: "EditAlbumDatesDialogOpened";
    albumId: AlbumId;
}

export function editAlbumDatesDialogOpened(albumId: AlbumId): EditAlbumDatesDialogOpened {
    return {
        type: "EditAlbumDatesDialogOpened",
        albumId,
    };
}

export function reduceEditAlbumDatesDialogOpened(
    current: CatalogViewerState,
    {albumId}: EditAlbumDatesDialogOpened,
): CatalogViewerState {
    return {
        ...current,
        editAlbumDatesDialog: {
            albumId,
        },
    };
}

export function editAlbumDatesDialogOpenedReducerRegistration(handlers: any) {
    handlers["EditAlbumDatesDialogOpened"] = reduceEditAlbumDatesDialogOpened as (
        state: CatalogViewerState,
        action: EditAlbumDatesDialogOpened
    ) => CatalogViewerState;
}
