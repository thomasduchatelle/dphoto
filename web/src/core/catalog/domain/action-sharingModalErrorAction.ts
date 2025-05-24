import {CatalogViewerState, ShareError} from "./catalog-state";

export interface SharingModalErrorAction {
    type: "SharingModalErrorAction"
    error: ShareError
}

export function sharingModalErrorAction(props: Omit<SharingModalErrorAction, "type">): SharingModalErrorAction {
    return {
        ...props,
        type: "SharingModalErrorAction",
    };
}

export function reduceSharingModalError(
    current: CatalogViewerState,
    action: SharingModalErrorAction,
): CatalogViewerState {
    if (!current.shareModal) return current;
    return {
        ...current,
        shareModal: {
            ...current.shareModal,
            error: action.error,
        }
    };
}
