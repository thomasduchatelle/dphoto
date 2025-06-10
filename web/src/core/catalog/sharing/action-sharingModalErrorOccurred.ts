import {CatalogViewerState, ShareError} from "../language";
import {moveSharedWithToSuggestion, moveSuggestionToSharedWith} from "./sharing";

export interface SharingModalErrorOccurred {
    type: "sharingModalErrorOccurred"
    error: ShareError
}

export function sharingModalErrorOccurred(error: ShareError): SharingModalErrorOccurred {
    return {
        error,
        type: "sharingModalErrorOccurred",
    };
}


export function reduceSharingModalErrorOccurred(
    current: CatalogViewerState,
    {error}: SharingModalErrorOccurred,
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

export function sharingModalErrorOccurredReducerRegistration(handlers: any) {
    handlers["sharingModalErrorOccurred"] = reduceSharingModalErrorOccurred as (
        state: CatalogViewerState,
        action: SharingModalErrorOccurred
    ) => CatalogViewerState;
}
