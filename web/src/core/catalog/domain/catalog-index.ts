import {OpenSharingModalAction, openSharingModalAction} from "./action-openSharingModalAction";
import {AddSharingAction, addSharingAction} from "./action-addSharingAction";
import {RemoveSharingAction, removeSharingAction} from "./action-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction} from "./action-closeSharingModalAction";
import {SharingModalErrorAction, sharingModalErrorAction} from "./action-sharingModalErrorAction";
import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction} from "./action-albumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction} from "./action-albumsLoadedAction";
import {MediaFailedToLoadAction, mediaFailedToLoadAction} from "./action-mediaFailedToLoadAction";
import {NoAlbumAvailableAction, noAlbumAvailableAction} from "./action-noAlbumAvailableAction";
import {StartLoadingMediasAction, startLoadingMediasAction} from "./action-startLoadingMediasAction";
import {AlbumsFilteredAction, albumsFilteredAction} from "./action-albumsFilteredAction";

export * from "./catalog-reducer";

export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
    sharingModalErrorAction,
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
    albumsFilteredAction,
};

export type {
    OpenSharingModalAction,
    AddSharingAction,
    RemoveSharingAction,
    CloseSharingModalAction,
    SharingModalErrorAction,
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
    MediaFailedToLoadAction,
    NoAlbumAvailableAction,
    StartLoadingMediasAction,
    AlbumsFilteredAction,
};
