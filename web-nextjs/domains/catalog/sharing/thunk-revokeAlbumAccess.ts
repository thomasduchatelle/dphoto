import {AlbumId, CatalogViewerState, isShareDialog} from "../language";
import {CatalogDispatch} from "../common/catalog-dispatch";
import {AlbumAccessRevoked, albumAccessRevoked} from "./action-albumAccessRevoked";
import {SharingModalErrorOccurred, sharingModalErrorOccurred} from "./action-sharingModalErrorOccurred";
import {ThunkDeclaration} from "@/libs/dthunks";

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
    CatalogDispatch & { adapter: RevokeAlbumAccessAPI }
> = {
    factory: ({dispatch, adapter, partialState: {albumId}}) => {
        return revokeAlbumAccessThunk.bind(null, dispatch, adapter, albumId);
    },
    selector: (state: CatalogViewerState) => {
        const dialog = state.dialog;
        if (!isShareDialog(dialog)) {
            return {albumId: undefined};
        }
        return {
            albumId: dialog.sharedAlbumId,
        };
    },
};
