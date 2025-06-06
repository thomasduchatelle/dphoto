import {CatalogViewerState, ShareError} from "../catalog-state";
import {moveSharedWithToSuggestion, moveSuggestionToSharedWith} from "./sharing";

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
    {error}: SharingModalErrorAction,
): CatalogViewerState {
    if (!current.shareModal) {
        return current;
    }

    if (error.type === "grant") {
        return moveSharedWithToSuggestion(current, current.shareModal, error.email, error);
    }
    if (error.type === "revoke") {
        const user = current.shareModal.suggestions.find(s => s.email === error.email) ?? {name: error.email, email: error.email}
        return moveSuggestionToSharedWith(current, current.shareModal, user, error);
    }

    return current;
}

export function sharingModalErrorReducerRegistration(handlers: any) {
    handlers["SharingModalErrorAction"] = reduceSharingModalError as (
        state: CatalogViewerState,
        action: SharingModalErrorAction
    ) => CatalogViewerState;
}
