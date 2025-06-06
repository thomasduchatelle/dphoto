import {AlbumId, catalogActions, CatalogViewerAction, CatalogViewerState} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "./catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";

export interface revokeAlbumSharingAPI {
    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void>;
}

export function revokeAlbumSharingThunk(
    dispatch: (action: CatalogViewerAction) => void,
    sharingAPI: revokeAlbumSharingAPI,
    albumId: AlbumId | undefined,
    email: string
): Promise<void> {
    if (!albumId) {
        return Promise.reject(`ERROR: no albumId selected to be revoked, cannot revoke access for ${email}`);
    }

    dispatch(catalogActions.removeSharingAction(email))

    return sharingAPI.revokeSharingAlbum(albumId, email)
        .catch(err => {
            console.log(`ERROR: ${JSON.stringify(err)}`);
            dispatch(catalogActions.sharingModalErrorAction({
                error: {type: "revoke", email, message: `Couldn't revoke access of user ${email}, try again later`}
            }));
        });
}

export const revokeAlbumSharingDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { albumId?: AlbumId },
    (email: string) => Promise<void>,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app, partialState: {albumId}}) => {
        const sharingAPI: revokeAlbumSharingAPI = new CatalogAPIAdapter(app.axiosInstance, app);
        return revokeAlbumSharingThunk.bind(null, dispatch, sharingAPI, albumId);
    },
    selector: ({shareModal}: CatalogViewerState) => ({
        albumId: shareModal?.sharedAlbumId,
    }),
};
