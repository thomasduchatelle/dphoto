import {Album, AlbumId, CatalogViewerState, isDeleteDialog} from "../language";

export interface DeleteDialogFrag {
    albums: Album[];
    initialSelectedAlbumId?: AlbumId;
    isOpen: boolean;
    isLoading: boolean;
    error?: string;
}

export function deleteDialogSelector({dialog}: CatalogViewerState): DeleteDialogFrag {
    if (!isDeleteDialog(dialog)) {
        return {
            albums: [],
            initialSelectedAlbumId: undefined,
            isOpen: false,
            isLoading: false,
            error: undefined,
        };
    }
    return {
        albums: dialog.deletableAlbums ?? [],
        initialSelectedAlbumId: dialog.initialSelectedAlbumId,
        isOpen: true,
        isLoading: dialog.isLoading ?? false,
        error: dialog.error,
    };
}
