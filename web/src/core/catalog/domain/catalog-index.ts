import {OpenSharingModalAction, openSharingModalAction} from "./action-openSharingModalAction";
import {AddSharingAction, addSharingAction} from "./action-addSharingAction";
import {RemoveSharingAction, removeSharingAction} from "./action-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction} from "./action-closeSharingModalAction";
import {SharingModalErrorAction, sharingModalErrorAction} from "./action-sharingModalErrorAction";
import {catalogReducer} from "./catalog-reducer";

export type CatalogSupportedActions =
    | OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction
    | CloseSharingModalAction
    | SharingModalErrorAction;

export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
    sharingModalErrorAction,
};

export {
    catalogReducer,
};

export type {
    OpenSharingModalAction,
    AddSharingAction,
    RemoveSharingAction,
    CloseSharingModalAction,
    SharingModalErrorAction,
};
