import {Album, AlbumId, SharingType, UserDetails} from "./catalog-state";
import {CatalogViewerAction, catalogActions} from "./catalog-reducer-v2";


export interface SharingAPI {
    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void>

    grantAccessToAlbum(albumId: AlbumId, email: string, role: SharingType): Promise<void>

    loadUserDetails(email: string): Promise<UserDetails>
}

export class ShareController {

    constructor(readonly dispatch: (action: CatalogViewerAction) => void,
                readonly sharingAPI: SharingAPI) {
    }

    public openSharingModal = (album: Album): void => {
        this.dispatch(catalogActions.openSharingModalAction({albumId: album.albumId}))
    }

    public onClose = (): void => {
        this.dispatch(catalogActions.closeSharingModalAction())
    }

    public revokeAccess = (albumId: AlbumId, email: string): Promise<void> => {
        if (!albumId) {
            return Promise.reject(`ERROR: no albumId selected to be revoked, cannot revoke access for ${email}`)
        }

        return this.sharingAPI.revokeSharingAlbum(albumId, email)
            .then(() => this.dispatch(catalogActions.removeSharingAction({email})))
            .catch(err => {
                console.log(`ERROR: ${JSON.stringify(err)}`)
                this.dispatch(catalogActions.sharingModalErrorAction({
                    error: {type: "general", message: `Couldn't revoke access of user ${email}, try again later`}
                }))
            })
    }

    public grantAccess = (albumId: AlbumId, email: string, role: SharingType): Promise<void> => {
        return this.sharingAPI.grantAccessToAlbum(albumId, email, role)
            .then(() =>
                this.sharingAPI.loadUserDetails(email)
                    .catch(err => {
                        console.log(`WARN: failed to load user details ${email}, ${JSON.stringify(err)}`)
                        return Promise.resolve({
                            email: email,
                            name: email,
                        } as UserDetails)
                    })
                    .then(details => {
                        this.dispatch(catalogActions.addSharingAction({
                            sharing: {
                                user: details,
                                role: role,
                            }
                        }))
                    })
            )
            .catch(err => {
                console.log(`ERROR: ${JSON.stringify(err)}`)
                this.dispatch(catalogActions.sharingModalErrorAction({
                    error: {
                        type: "adding",
                        message: "Failed to grant access, verify the email address or contact maintainers"
                    }
                }))
            })
    }
}
