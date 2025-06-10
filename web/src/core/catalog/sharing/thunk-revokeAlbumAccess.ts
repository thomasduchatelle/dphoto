import {AlbumId, CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {AlbumAccessRevoked, albumAccessRevoked} from "./action-albumAccessRevoked";
import {SharingModalErrorOccurred, sharingModalErrorOccurred} from "./action-sharingModalErrorOccurred";

export interface RevokeAlbumAccessAPI {
    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void>;
}

export function revokeAlbumAccessThunk(
    dispatch: (action: AlbumAccessRevoked | SharingModalErrorOccurred) => void,
    sharingAPI: RevokeAlbumAccessAPI,
    albumId: AlbumId | undefined,
    email: string
): Promise<void> {
    if (!albumId) {
        return Promise.reject(`ERROR: no albumId selected to be revoked, cannot revoke access for ${email}`);
    }

    dispatch(albumAccessRevoked(email))

    return sharingAPI.revokeSharingAlbum(albumId, email)
        .catch(err => {
            console.log(`ERROR: ${JSON.stringify(err)}`);
            dispatch(sharingModalErrorOccurred({
                type: "revoke", email, message: `Couldn't revoke access of user ${email}, try again later`
            }));
        });
}

export const revokeAlbumAccessDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { albumId?: AlbumId },
    (email: string) => Promise<void>,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app, partialState: {albumId}}) => {
        const sharingAPI: RevokeAlbumAccessAPI = new CatalogAPIAdapter(app.axiosInstance, app);
        return revokeAlbumAccessThunk.bind(null, dispatch, sharingAPI, albumId);
    },
    selector: ({shareModal}: CatalogViewerState) => ({
        albumId: shareModal?.sharedAlbumId,
    }),
};
