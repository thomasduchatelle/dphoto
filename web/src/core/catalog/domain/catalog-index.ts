import {OpenSharingModalAction, openSharingModalAction} from "./action-openSharingModalAction";
import {AddSharingAction, addSharingAction} from "./action-addSharingAction";
import {RemoveSharingAction, removeSharingAction} from "./action-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction} from "./action-closeSharingModalAction";
import {catalogReducer} from "./catalog-reducer";

export type CatalogSupportedActions =
    | OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction
    | CloseSharingModalAction;

export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
};

export {
    catalogReducer,
};

export type {
    OpenSharingModalAction,
    AddSharingAction,
    RemoveSharingAction,
    CloseSharingModalAction,
};
