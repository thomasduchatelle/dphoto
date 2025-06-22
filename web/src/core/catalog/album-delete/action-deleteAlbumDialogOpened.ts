import {Album, AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

function isDeletable(album: Album): boolean {
    return albumIsOwnedByCurrentUser(album);
}

function getInitialSelectedAlbumId(
    deletableAlbums: Album[],
    loadingMediasFor?: AlbumId,
    mediasLoadedFromAlbumId?: AlbumId
): AlbumId | undefined {
    const albumId = loadingMediasFor ?? mediasLoadedFromAlbumId

    if (deletableAlbums && (!albumId || !deletableAlbums.some(a => isAlbumIdEqual(a.albumId, albumId)))) {
        return deletableAlbums[0]?.albumId; // Added optional chaining for safety
    }

    return albumId;
}

function isAlbumIdEqual(a: AlbumId, b: AlbumId): boolean {
    return a.owner === b.owner && a.folderName === b.folderName;
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
        return {
            ...current,
            deleteDialog: {
                deletableAlbums,
                initialSelectedAlbumId,
                isLoading: false,
                error: undefined,
            }
        };
    }
);

export type DeleteAlbumDialogOpened = ReturnType<typeof deleteAlbumDialogOpened>;
