import {CatalogViewerState} from "../language";
import {moveSharedWithToSuggestion} from "./sharing";
import {createAction} from "src/libs/daction";

export const albumAccessRevoked = createAction<CatalogViewerState, string>(
    "albumAccessRevoked",
    (current: CatalogViewerState, email: string) => {
        if (!current.shareModal) {
            return current;
        }

        return moveSharedWithToSuggestion(current, current.shareModal, email);
    }
);

export type AlbumAccessRevoked = ReturnType<typeof albumAccessRevoked>;
