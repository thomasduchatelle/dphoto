import {OpenSharingModalAction, openSharingModalAction} from "./action-openSharingModalAction";
import {AddSharingAction, addSharingAction} from "./action-addSharingAction";
import {RemoveSharingAction, removeSharingAction} from "./action-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction} from "./action-closeSharingModalAction";
import {SharingModalErrorAction, sharingModalErrorAction} from "./action-sharingModalErrorAction";

export * from "./catalog-reducer";


export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
    sharingModalErrorAction,
};

export type {
    OpenSharingModalAction,
    AddSharingAction,
    RemoveSharingAction,
    CloseSharingModalAction,
    SharingModalErrorAction,
};
