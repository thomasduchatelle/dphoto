import {CatalogViewerState} from "./catalog-state";
import {reduceAlbumsAndMediasLoaded} from "./action-albumsAndMediasLoadedAction";
import {reduceAlbumsLoaded} from "./action-albumsLoadedAction";
import {reduceMediaFailedToLoad} from "./action-mediaFailedToLoadAction";
import {reduceNoAlbumAvailable} from "./action-noAlbumAvailableAction";
import {reduceStartLoadingMedias} from "./action-startLoadingMediasAction";
import {reduceAlbumsFiltered} from "./action-albumsFilteredAction";
import {reduceOpenSharingModal} from "./action-openSharingModalAction";
import {reduceAddSharing} from "./action-addSharingAction";
import {reduceRemoveSharing} from "./action-removeSharingAction";
import {CatalogViewerAction} from "./catalog-actions";
import {reduceCloseSharingModal} from "./action-closeSharingModalAction";
import {reduceSharingModalError} from "./action-sharingModalErrorAction";

export function catalogReducer(
    state: CatalogViewerState,
    action: CatalogViewerAction
): CatalogViewerState {
    switch (action.type) {
        case "AlbumsAndMediasLoadedAction":
            return reduceAlbumsAndMediasLoaded(state, action);
        case "AlbumsLoadedAction":
            return reduceAlbumsLoaded(state, action);
        case "MediaFailedToLoadAction":
            return reduceMediaFailedToLoad(state, action);
        case "NoAlbumAvailableAction":
            return reduceNoAlbumAvailable(state, action);
        case "StartLoadingMediasAction":
            return reduceStartLoadingMedias(state, action);
        case "AlbumsFilteredAction":
            return reduceAlbumsFiltered(state, action);
        case "OpenSharingModalAction":
            return reduceOpenSharingModal(state, action);
        case "AddSharingAction":
            return reduceAddSharing(state, action);
        case "RemoveSharingAction":
            return reduceRemoveSharing(state, action);
        case "CloseSharingModalAction":
            return reduceCloseSharingModal(state, action);
        case "SharingModalErrorAction":
            return reduceSharingModalError(state, action);
        default:
            return state;
    }
}
