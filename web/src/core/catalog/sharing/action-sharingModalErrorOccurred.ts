import {CatalogViewerState, ShareError, isShareDialog} from "../language";
import {moveSharedWithToSuggestion, moveSuggestionToSharedWith} from "./sharing";
import {createAction} from "src/libs/daction";

export const sharingModalErrorOccurred = createAction<CatalogViewerState, ShareError>(
    "sharingModalErrorOccurred",
    (current: CatalogViewerState, error: ShareError) => {
        if (!isShareDialog(current.dialog)) {
            return current;
        }

        if (error.type === "grant") {
            return moveSharedWithToSuggestion(current, current.dialog, error.email, error);
        }
        if (error.type === "revoke") {
            const user = current.dialog.suggestions.find(s => s.email === error.email) ?? {name: error.email, email: error.email}
            return moveSuggestionToSharedWith(current, current.dialog, user, error);
        }

        return current;
    }
);

export type SharingModalErrorOccurred = ReturnType<typeof sharingModalErrorOccurred>;
