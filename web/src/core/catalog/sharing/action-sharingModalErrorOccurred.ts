import {CatalogViewerState, ShareError} from "../language";
import {moveSharedWithToSuggestion, moveSuggestionToSharedWith} from "./sharing";
import {createAction} from "src/light-state-lib";

export const sharingModalErrorOccurred = createAction<CatalogViewerState, ShareError>(
    "sharingModalErrorOccurred",
    (current: CatalogViewerState, error: ShareError) => {
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
);

export type SharingModalErrorOccurred = ReturnType<typeof sharingModalErrorOccurred>;
