import {Album, AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "../catalog-state";

export interface OpenDeleteAlbumDialogAction {
    type: "OpenDeleteAlbumDialog"
}

export function openDeleteAlbumDialogAction(): OpenDeleteAlbumDialogAction {
    return {
        type: "OpenDeleteAlbumDialog",
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

export function reduceOpenDeleteAlbumDialog(
    current: CatalogViewerState,
    action: OpenDeleteAlbumDialogAction,
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

export function openDeleteAlbumDialogReducerRegistration(handlers: any) {
    handlers["OpenDeleteAlbumDialog"] = reduceOpenDeleteAlbumDialog as (
        state: CatalogViewerState,
        action: OpenDeleteAlbumDialogAction
    ) => CatalogViewerState;
}
