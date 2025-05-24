import {CatalogViewerState} from "./catalog-state";

export interface CloseSharingModalAction {
    type: "CloseSharingModalAction"
}

export function closeSharingModalAction(): CloseSharingModalAction {
    return {type: "CloseSharingModalAction"};
}

export function reduceCloseSharingModal(
    current: CatalogViewerState,
    action: CloseSharingModalAction,
): CatalogViewerState {
    const {shareModal, ...rest} = current;
    return rest as CatalogViewerState;
}

export function closeSharingModalReducerRegistration(handlers: any) {
    handlers["CloseSharingModalAction"] = reduceCloseSharingModal as (
        state: CatalogViewerState,
        action: CloseSharingModalAction
    ) => CatalogViewerState;
}
