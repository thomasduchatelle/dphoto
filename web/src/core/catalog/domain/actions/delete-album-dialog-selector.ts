import {Album, AlbumId, CatalogViewerState} from "../catalog-state";

export interface DeleteAlbumDialogSelectorResult {
    albums: Album[];
    initialSelectedAlbumId?: AlbumId;
    isOpen: boolean;
    isLoading: boolean;
    error?: string;
}

export function selectDeleteAlbumDialog({deleteDialog}: CatalogViewerState): DeleteAlbumDialogSelectorResult {
    return {
        albums: deleteDialog?.deletableAlbums ?? [],
        initialSelectedAlbumId: deleteDialog?.initialSelectedAlbumId,
        isOpen: !!deleteDialog,
        isLoading: deleteDialog?.isLoading ?? false,
        error: deleteDialog?.error,
    };
}
