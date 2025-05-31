import {CatalogViewerState, ShareError} from "../catalog-state";

export interface SharingModalErrorAction {
    type: "SharingModalErrorAction"
    error: ShareError
}

export function sharingModalErrorAction(props: ShareError | Omit<SharingModalErrorAction, "type">): SharingModalErrorAction {
    if ("type" in props && "message" in props) {
        return {
            type: "SharingModalErrorAction",
            error: props as ShareError,
        }
    }
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

export function sharingModalErrorReducerRegistration(handlers: any) {
    handlers["SharingModalErrorAction"] = reduceSharingModalError as (
        state: CatalogViewerState,
        action: SharingModalErrorAction
    ) => CatalogViewerState;
}
