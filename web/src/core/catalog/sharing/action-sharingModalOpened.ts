import {AlbumId, CatalogViewerState} from "../language";
import {withOpenShareModal} from "./sharing";

export interface SharingModalOpened {
    type: "sharingModalOpened"
    albumId: AlbumId
}

export function sharingModalOpened(albumId: AlbumId): SharingModalOpened {
    return {
        albumId,
        type: "sharingModalOpened",
    };
}

export function reduceSharingModalOpened(
    current: CatalogViewerState,
    action: SharingModalOpened,
): CatalogViewerState {
    return withOpenShareModal(current, action.albumId);
}

export function sharingModalOpenedReducerRegistration(handlers: any) {
    handlers["sharingModalOpened"] = reduceSharingModalOpened as (
        state: CatalogViewerState,
        action: SharingModalOpened
    ) => CatalogViewerState;
}
