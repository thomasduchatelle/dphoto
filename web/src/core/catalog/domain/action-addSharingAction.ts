import {CatalogViewerState, Sharing} from "./catalog-state";

export interface AddSharingAction {
    type: "AddSharingAction"
    sharing: Sharing
}

export function addSharingAction(props: Omit<AddSharingAction, "type">): AddSharingAction {
    return {
        ...props,
        type: "AddSharingAction",
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
