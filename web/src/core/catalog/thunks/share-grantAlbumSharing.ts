import {AlbumId, catalogActions, CatalogViewerAction, CatalogViewerState, UserDetails} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "./catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";

export interface GrantAlbumSharingAPI {
    grantAccessToAlbum(albumId: AlbumId, email: string): Promise<void>;

    loadUserDetails(email: string): Promise<UserDetails>;
}

export function grantAlbumSharingThunk(
    dispatch: (action: CatalogViewerAction) => void,
    sharingAPI: GrantAlbumSharingAPI,
    albumId: AlbumId | undefined,
    email: string
): Promise<void> {
    if (!albumId) {
        return Promise.reject(`ERROR: no albumId selected to be granted, cannot grant access for ${email}`);
    }
    return sharingAPI.grantAccessToAlbum(albumId, email)
        .then(() =>
            sharingAPI.loadUserDetails(email)
                .catch(err => {
                    console.log(`WARN: failed to load user details ${email}, ${JSON.stringify(err)}`);
                    return Promise.resolve({
                        email: email,
                        name: email,
                    } as UserDetails);
                })
                .then(details => {
                    dispatch(catalogActions.addSharingAction({
                        sharing: {
                            user: details,
                        }
                    }));
                })
        )
        .catch(err => {
            console.log(`ERROR: ${JSON.stringify(err)}`);
            dispatch(catalogActions.sharingModalErrorAction({
                error: {
                    type: "adding",
                    message: "Failed to grant access, verify the email address or contact maintainers"
                }
            }));
        });
}

export const grantAlbumSharingDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { albumId?: AlbumId },
    (email: string) => Promise<void>,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app, partialState: {albumId}}) => {
        const sharingAPI: GrantAlbumSharingAPI = new CatalogAPIAdapter(app.axiosInstance, app);
        return grantAlbumSharingThunk.bind(null, dispatch, sharingAPI, albumId);
    },
    selector: ({shareModal}: CatalogViewerState) => ({
        albumId: shareModal?.sharedAlbumId,
    }),
};
