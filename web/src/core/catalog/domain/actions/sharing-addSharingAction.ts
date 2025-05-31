import {CatalogViewerState, Sharing} from "../catalog-state";
import {sortSharings} from "./sharing-openSharingModalAction";

export interface AddSharingAction {
    type: "AddSharingAction"
    sharing: Sharing
}

export function addSharingAction(props: Sharing | Omit<AddSharingAction, "type">): AddSharingAction {
    if ("user" in props && "role" in props) {
        return {
            type: "AddSharingAction",
            sharing: props,
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
    return {
        ...current,
        shareModal: {
            ...current.shareModal,
            sharedWith: sortSharings(updatedSharedWith),
        }
    };
}

export function addSharingReducerRegistration(handlers: any) {
    handlers["AddSharingAction"] = reduceAddSharing as (state: CatalogViewerState, action: AddSharingAction) => CatalogViewerState;
}
