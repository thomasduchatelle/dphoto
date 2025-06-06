import {AlbumId, CatalogViewerState} from "../catalog-state";
import {withOpenShareModal} from "./sharing";

export interface OpenSharingModalAction {
    type: "OpenSharingModalAction"
    albumId: AlbumId
}

export function openSharingModalAction(props: AlbumId | Omit<OpenSharingModalAction, "type">): OpenSharingModalAction {
    if ("owner" in props && "folderName" in props) {
        return {
            type: "OpenSharingModalAction",
            albumId: props,
        };
    }
    return {
        ...props,
        type: "OpenSharingModalAction",
    };
}

export function reduceOpenSharingModal(
    current: CatalogViewerState,
    action: OpenSharingModalAction,
): CatalogViewerState {
    return withOpenShareModal(current, action.albumId);
}

export function openSharingModalReducerRegistration(handlers: any) {
    handlers["OpenSharingModalAction"] = reduceOpenSharingModal as (
        state: CatalogViewerState,
        action: OpenSharingModalAction
    ) => CatalogViewerState;
}
