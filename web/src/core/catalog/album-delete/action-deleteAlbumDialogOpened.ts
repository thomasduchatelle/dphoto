import {Album, AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "../language";

export interface DeleteAlbumDialogOpened {
    type: "deleteAlbumDialogOpened"
}

export function deleteAlbumDialogOpened(): DeleteAlbumDialogOpened {
    return {
        type: "deleteAlbumDialogOpened",
    };
}

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
        return deletableAlbums[0].albumId;
    }

    return albumId;
}

function isAlbumIdEqual(a: AlbumId, b: AlbumId): boolean {
    return a.owner === b.owner && a.folderName === b.folderName;
}

export function reduceDeleteAlbumDialogOpened(
    current: CatalogViewerState,
    action: DeleteAlbumDialogOpened,
): CatalogViewerState {
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

export function deleteAlbumDialogOpenedReducerRegistration(handlers: any) {
    handlers["deleteAlbumDialogOpened"] = reduceDeleteAlbumDialogOpened as (
        state: CatalogViewerState,
        action: DeleteAlbumDialogOpened
    ) => CatalogViewerState;
}
