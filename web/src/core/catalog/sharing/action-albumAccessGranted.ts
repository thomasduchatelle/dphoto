import {CatalogViewerState, Sharing} from "../language";
import {moveSuggestionToSharedWith} from "./sharing";

export interface AlbumAccessGranted {
    type: "albumAccessGranted"
    sharing: Sharing
}

export function albumAccessGranted(sharing: Sharing): AlbumAccessGranted {
    return {
        sharing,
        type: "albumAccessGranted",
    };
}

export function reduceAlbumAccessGranted(
    current: CatalogViewerState,
    action: AlbumAccessGranted,
): CatalogViewerState {
    if (!current.shareModal) return current;

    return moveSuggestionToSharedWith(current, current.shareModal, action.sharing.user);
}

export function albumAccessGrantedReducerRegistration(handlers: any) {
    handlers["albumAccessGranted"] = reduceAlbumAccessGranted as (state: CatalogViewerState, action: AlbumAccessGranted) => CatalogViewerState;
}
