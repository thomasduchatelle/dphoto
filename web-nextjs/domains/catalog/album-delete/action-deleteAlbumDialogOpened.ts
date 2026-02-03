import {Album, AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState, DeleteDialog} from "../language";
import {createAction} from "@/libs/daction";
import {albumIdEquals} from "../language/utils-albumIdEquals";

function isDeletable(album: Album): boolean {
    return albumIsOwnedByCurrentUser(album);
}

function getInitialSelectedAlbumId(
    deletableAlbums: Album[],
    loadingMediasFor?: AlbumId,
    mediasLoadedFromAlbumId?: AlbumId
): AlbumId | undefined {
    const albumId = loadingMediasFor ?? mediasLoadedFromAlbumId

    if (deletableAlbums && (!albumId || !deletableAlbums.some(a => albumIdEquals(a.albumId, albumId)))) {
        return deletableAlbums[0]?.albumId; // Added optional chaining for safety
    }

    return albumId;
}

export const deleteAlbumDialogOpened = createAction<CatalogViewerState>(
    "deleteAlbumDialogOpened",
    (current: CatalogViewerState) => {
        const deletableAlbums = (current.allAlbums ?? []).filter(isDeletable);
        const initialSelectedAlbumId = getInitialSelectedAlbumId(
            deletableAlbums,
            current.loadingMediasFor,
            current.mediasLoadedFromAlbumId
        );
        const deleteDialog: DeleteDialog = {
            type: "DeleteDialog",
            deletableAlbums,
            initialSelectedAlbumId,
            isLoading: false,
            error: undefined,
        }
        return {
            ...current,
            dialog: deleteDialog
        };
    }
);

export type DeleteAlbumDialogOpened = ReturnType<typeof deleteAlbumDialogOpened>;
