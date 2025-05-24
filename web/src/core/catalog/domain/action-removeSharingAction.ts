import {CatalogViewerState} from "./catalog-state";

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
    return {
        ...current,
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
