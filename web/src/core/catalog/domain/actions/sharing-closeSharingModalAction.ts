import {CatalogViewerState} from "../catalog-state";

export interface CloseSharingModalAction {
    type: "CloseSharingModalAction"
}

export function closeSharingModalAction(): CloseSharingModalAction {
    return {type: "CloseSharingModalAction"};
}

export function reduceCloseSharingModal(
    {shareModal, ...rest}: CatalogViewerState,
    _: CloseSharingModalAction,
): CatalogViewerState {
    return rest;
}

export function closeSharingModalReducerRegistration(handlers: any) {
    handlers["CloseSharingModalAction"] = reduceCloseSharingModal as (
        state: CatalogViewerState,
        action: CloseSharingModalAction
    ) => CatalogViewerState;
}
