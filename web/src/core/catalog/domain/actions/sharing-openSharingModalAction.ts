import {AlbumId, CatalogViewerState, Sharing} from "../catalog-state";
import {albumIdEquals} from "../utils-albumIdEquals";

export interface OpenSharingModalAction {
    type: "OpenSharingModalAction"
    albumId: AlbumId
}

export function openSharingModalAction(props: AlbumId | Omit<OpenSharingModalAction, "type">): OpenSharingModalAction {
    if ("owner" in props && "folderName" in props) {
        return {
            type: "OpenSharingModalAction",
            albumId: props,
        };
    }
    return {
        ...props,
        type: "OpenSharingModalAction",
    };
}

export function sortSharings(sharings: Sharing[]): Sharing[] {
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
    const album = current.allAlbums.find(a => albumIdEquals(a.albumId, action.albumId));
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

export function openSharingModalReducerRegistration(handlers: any) {
    handlers["OpenSharingModalAction"] = reduceOpenSharingModal as (
        state: CatalogViewerState,
        action: OpenSharingModalAction
    ) => CatalogViewerState;
}
