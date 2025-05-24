import {AlbumId, CatalogViewerState, Sharing} from "./catalog-state";

export interface OpenSharingModalAction {
    type: "OpenSharingModalAction"
    albumId: AlbumId
}

export function openSharingModalAction(props: Omit<OpenSharingModalAction, "type">): OpenSharingModalAction {
    return {
        ...props,
        type: "OpenSharingModalAction",
    };
}

function sortSharings(sharings: Sharing[]): Sharing[] {
    return sharings.slice().sort((a, b) => {
        const nameA = a.user.name?.trim() || "";
        const nameB = b.user.name?.trim() || "";
        if (nameA && nameB) {
            const cmp = nameA.localeCompare(nameB);
            if (cmp !== 0) return cmp;
            return a.user.email.localeCompare(b.user.email);
        }
        if (!nameA && !nameB) {
            return a.user.email.localeCompare(b.user.email);
        }
        if (!nameA) return 1;
        if (!nameB) return -1;
        return 0;
    });
}

export function reduceOpenSharingModal(
    current: CatalogViewerState,
    action: OpenSharingModalAction,
): CatalogViewerState {
    // Find the album in the current state
    const album = current.allAlbums.find(a => a.albumId.owner === action.albumId.owner && a.albumId.folderName === action.albumId.folderName);
    return {
        ...current,
        shareModal: album
            ? {
                sharedAlbumId: album.albumId,
                sharedWith: sortSharings([...album.sharedWith]),
            }
            : undefined,
    };
}
