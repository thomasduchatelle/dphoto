import {AlbumId, CatalogViewerState, UserDetails} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {AlbumAccessGranted, albumAccessGranted} from "./action-albumAccessGranted";
import {SharingModalErrorOccurred, sharingModalErrorOccurred} from "./action-sharingModalErrorOccurred";

export interface GrantAlbumAccessAPI {
    grantAccessToAlbum(albumId: AlbumId, email: string): Promise<void>;

    loadUserDetails(email: string): Promise<UserDetails>;
}

export function grantAlbumAccessThunk(
    dispatch: (action: AlbumAccessGranted | SharingModalErrorOccurred) => void,
    sharingAPI: GrantAlbumAccessAPI,
    albumId: AlbumId | undefined,
    email: string
): Promise<void> {
    if (!albumId) {
        return Promise.reject(`ERROR: no albumId selected to be granted, cannot grant access for ${email}`);
    }

    return Promise.allSettled([
        sharingAPI.grantAccessToAlbum(albumId, email),
        sharingAPI.loadUserDetails(email)
            .catch(err => {
                console.log(`WARN: failed to load user details ${email}, ${JSON.stringify(err)}`);
                return Promise.resolve({
                    email: email,
                    name: email,
                } as UserDetails);
            })
            .then(details => {
                dispatch(albumAccessGranted({
                    user: details,
                }));
            })
    ])
        .then(([grantResp, _]) => {
            if (grantResp.status === "rejected") {
                console.log(`ERROR: ${JSON.stringify(grantResp.reason)}`);
                dispatch(sharingModalErrorOccurred({
                    type: "grant",
                    message: "Failed to grant access, verify the email address or contact maintainers",
                    email: email,
                }));
                return Promise.reject(grantResp.reason);
            }
        })
}

export const grantAlbumAccessDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { albumId?: AlbumId },
    (email: string) => Promise<void>,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app, partialState: {albumId}}) => {
        const sharingAPI: GrantAlbumAccessAPI = new CatalogAPIAdapter(app.axiosInstance, app);
        return grantAlbumAccessThunk.bind(null, dispatch, sharingAPI, albumId);
    },
    selector: ({shareModal}: CatalogViewerState) => ({
        albumId: shareModal?.sharedAlbumId,
    }),
};
