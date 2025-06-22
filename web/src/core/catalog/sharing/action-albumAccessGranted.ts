import {CatalogViewerState, Sharing} from "../language";
import {moveSuggestionToSharedWith} from "./sharing";
import {createAction} from "src/light-state-lib";

export const albumAccessGranted = createAction<CatalogViewerState, Sharing>(
    "albumAccessGranted",
    (current: CatalogViewerState, sharing: Sharing) => {
        if (!current.shareModal) return current;

        return moveSuggestionToSharedWith(current, current.shareModal, sharing.user);
    }
);

export type AlbumAccessGranted = ReturnType<typeof albumAccessGranted>;
