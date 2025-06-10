import {CatalogViewerState} from "../language";

import {moveSharedWithToSuggestion} from "./sharing";

export interface AlbumAccessRevoked {
    type: "albumAccessRevoked"
    email: string
}

export function albumAccessRevoked(email: string): AlbumAccessRevoked {
    return {
        email,
        type: "albumAccessRevoked",
    };
}

export function reduceAlbumAccessRevoked(
    current: CatalogViewerState,
    action: AlbumAccessRevoked,
): CatalogViewerState {
    if (!current.shareModal) {
        return current;
    }

    return moveSharedWithToSuggestion(current, current.shareModal, action.email);
}

export function albumAccessRevokedReducerRegistration(handlers: any) {
    handlers["albumAccessRevoked"] = reduceAlbumAccessRevoked as (
        state: CatalogViewerState,
        action: AlbumAccessRevoked
    ) => CatalogViewerState;
}
