import {CatalogViewerState} from "../language";

export interface SharingModalClosed {
    type: "sharingModalClosed";
}

export function sharingModalClosed(): SharingModalClosed {
    return {type: "sharingModalClosed"};
}

export function reduceSharingModalClosed(
    {shareModal, ...rest}: CatalogViewerState,
    _: SharingModalClosed,
): CatalogViewerState {
    return rest;
}

export function sharingModalClosedReducerRegistration(handlers: any) {
    handlers["sharingModalClosed"] = reduceSharingModalClosed as (
        state: CatalogViewerState,
        action: SharingModalClosed
    ) => CatalogViewerState;
}
