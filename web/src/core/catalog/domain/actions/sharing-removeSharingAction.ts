import {albumMatchCriterion, CatalogViewerState} from "../catalog-state";
import {albumIdEquals} from "../utils-albumIdEquals";

export interface RemoveSharingAction {
    type: "RemoveSharingAction"
    email: string
}

export function removeSharingAction(props: string | Omit<RemoveSharingAction, "type">): RemoveSharingAction {
    if (typeof props === "string") {
        return {
            type: "RemoveSharingAction",
            email: props,
        };
    }
    return {
        ...props,
        type: "RemoveSharingAction",
    };
}

export function reduceRemoveSharing(
    current: CatalogViewerState,
    action: RemoveSharingAction,
): CatalogViewerState {
    if (!current.shareModal) return current;
    const updatedSharedWith = current.shareModal.sharedWith.filter(
        s => s.user.email !== action.email
    );

    const updatedAllAlbums = current.allAlbums.map(album => {
        if (current.shareModal && albumIdEquals(album.albumId, current.shareModal.sharedAlbumId)) {
            const albumUpdatedSharedWith = album.sharedWith.filter(s => s.user.email !== action.email);
            return {
                ...album,
                sharedWith: albumUpdatedSharedWith,
            };
        }
        return album;
    });

    return {
        ...current,
        albums: updatedAllAlbums.filter(albumMatchCriterion(current.albumFilter.criterion)),
        allAlbums: updatedAllAlbums,
        shareModal: {
            ...current.shareModal,
            sharedWith: updatedSharedWith,
        }
    };
}

export function removeSharingReducerRegistration(handlers: any) {
    handlers["RemoveSharingAction"] = reduceRemoveSharing as (
        state: CatalogViewerState,
        action: RemoveSharingAction
    ) => CatalogViewerState;
}
