import {CatalogViewerState, Sharing} from "../catalog-state";
import {moveSuggestionToSharedWith} from "./sharing";

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

    return moveSuggestionToSharedWith(current, current.shareModal, action.sharing.user);
    // // Replace if user already exists (by email), else add
    // const newSharing = action.sharing;
    // const updatedSharedWith = [
    //     ...current.shareModal.sharedWith.filter(s => s.user.email !== newSharing.user.email),
    //     newSharing
    // ];
    //
    // // Also update allAlbums for consistency
    // const updatedAllAlbums = current.allAlbums.map(album => {
    //     if (current.shareModal && albumIdEquals(album.albumId, current.shareModal.sharedAlbumId)) {
    //         return {
    //             ...album,
    //             sharedWith: sortSharings(updatedSharedWith),
    //         };
    //     }
    //     return album;
    // });
    //
    // // Remove the newly granted email from suggestions if present
    // const grantedEmail = newSharing.user.email;
    // const updatedSuggestions = current.shareModal.suggestions
    //     ? current.shareModal.suggestions.filter(s => s.email !== grantedEmail)
    //     : [];
    //
    // return {
    //     ...current,
    //     albums: updatedAllAlbums.filter(albumMatchCriterion(current.albumFilter.criterion)),
    //     allAlbums: updatedAllAlbums,
    //     shareModal: {
    //         ...current.shareModal,
    //         sharedWith: sortSharings(updatedSharedWith),
    //         suggestions: updatedSuggestions,
    //     }
    // };
}

export function addSharingReducerRegistration(handlers: any) {
    handlers["AddSharingAction"] = reduceAddSharing as (state: CatalogViewerState, action: AddSharingAction) => CatalogViewerState;
}
