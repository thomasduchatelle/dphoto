import {CatalogViewerState, isShareDialog} from "../language";
import {moveSharedWithToSuggestion} from "./sharing";
import {createAction} from "@/libs/daction";

export const albumAccessRevoked = createAction<CatalogViewerState, string>(
    "albumAccessRevoked",
    (current: CatalogViewerState, email: string) => {
        if (!isShareDialog(current.dialog)) {
            return current;
        }

        return moveSharedWithToSuggestion(current, current.dialog, email);
    }
);

export type AlbumAccessRevoked = ReturnType<typeof albumAccessRevoked>;
