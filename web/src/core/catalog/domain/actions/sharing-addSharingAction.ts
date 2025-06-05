import {albumMatchCriterion, CatalogViewerState, Sharing} from "../catalog-state";
import {sortSharings} from "./sharing-openSharingModalAction";
import {albumIdEquals} from "../utils-albumIdEquals";

export interface AddSharingAction {
    type: "AddSharingAction"
    sharing: Sharing
}

export function addSharingAction(props: Sharing | Omit<AddSharingAction, "type">): AddSharingAction {
    if ("user" in props) {
        return {
            type: "AddSharingAction",
            sharing: {user: (props as Sharing).user},
        };
    }
    return {
        ...props,
        type: "AddSharingAction",
    };
}

export function reduceAddSharing(
    current: CatalogViewerState,
    action: AddSharingAction,
): CatalogViewerState {
    if (!current.shareModal) return current;
    // Replace if user already exists (by email), else add
    const newSharing = action.sharing;
    const updatedSharedWith = [
        ...current.shareModal.sharedWith.filter(s => s.user.email !== newSharing.user.email),
        newSharing
    ];

    // Also update allAlbums for consistency
    const updatedAllAlbums = current.allAlbums.map(album => {
        if (current.shareModal && albumIdEquals(album.albumId, current.shareModal.sharedAlbumId)) {
            return {
                ...album,
                sharedWith: sortSharings(updatedSharedWith),
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
            sharedWith: sortSharings(updatedSharedWith),
        }
    };
}

export function addSharingReducerRegistration(handlers: any) {
    handlers["AddSharingAction"] = reduceAddSharing as (state: CatalogViewerState, action: AddSharingAction) => CatalogViewerState;
}
