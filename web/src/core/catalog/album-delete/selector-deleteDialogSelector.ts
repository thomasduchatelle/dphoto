import {Album, AlbumId, CatalogViewerState} from "../language";

export interface DeleteDialogFrag {
    albums: Album[];
    initialSelectedAlbumId?: AlbumId;
    isOpen: boolean;
    isLoading: boolean;
    error?: string;
}

export function deleteDialogSelector({deleteDialog}: CatalogViewerState): DeleteDialogFrag {
    return {
        albums: deleteDialog?.deletableAlbums ?? [],
        initialSelectedAlbumId: deleteDialog?.initialSelectedAlbumId,
        isOpen: !!deleteDialog,
        isLoading: deleteDialog?.isLoading ?? false,
        error: deleteDialog?.error,
    };
}
