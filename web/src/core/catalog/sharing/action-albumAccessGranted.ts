import {CatalogViewerState, Sharing, isShareDialog} from "../language";
import {moveSuggestionToSharedWith} from "./sharing";
import {createAction} from "src/libs/daction";

export const albumAccessGranted = createAction<CatalogViewerState, Sharing>(
    "albumAccessGranted",
    (current: CatalogViewerState, sharing: Sharing) => {
        if (!isShareDialog(current.dialog)) return current;

        return moveSuggestionToSharedWith(current, current.dialog, sharing.user);
    }
);

export type AlbumAccessGranted = ReturnType<typeof albumAccessGranted>;
