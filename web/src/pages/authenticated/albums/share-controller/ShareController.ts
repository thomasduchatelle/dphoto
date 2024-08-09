import {Album, AlbumId, SharingType, UserDetails} from "../../../../core/catalog-react";

import {SharingModalAction} from "./sharingModalReducer";

export interface SharingAPI {
    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void>;

    grantAccessToAlbum(albumId: AlbumId, email: string, role: SharingType): Promise<void>;

    loadUserDetails(email: string): Promise<UserDetails>
}

export class ShareController {

    constructor(readonly dispatch: (action: SharingModalAction) => void,
                readonly sharingAPI: SharingAPI,
                public albumId: AlbumId | undefined = undefined) {
    }

    public openSharingModal = (album: Album): void => {
        this.albumId = album.albumId
        this.dispatch({type: "OpenSharingModalAction", sharedWith: album.sharedWith})
    }

    public onClose = (): void => {
        this.albumId = undefined
        this.dispatch({type: "CloseSharingModalAction"})
    }

    public onRevoke = (email: string): void => {
        if (this.albumId) {
            this.sharingAPI.revokeSharingAlbum(this.albumId, email)
                .then(() => this.dispatch({type: "RemoveSharingAction", email}))
                .catch(err => {
                    console.log(`ERROR: ${JSON.stringify(err)}`)
                    this.dispatch({
                        type: "SharingModalErrorAction",
                        error: {type: "general", message: `Couldn't revoke access of user ${email}, try again later`}
                    })
                })
        }
    }

    public onGrant = (email: string, role: SharingType): Promise<void> => {
        if (this.albumId) {
            return this.sharingAPI.grantAccessToAlbum(this.albumId, email, role)
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
                            this.dispatch({
                                type: "AddSharingAction", sharing: {
                                    user: details,
                                    role: role,
                                }
                            })
                        })
                )
                .catch(err => {
                    console.log(`ERROR: ${JSON.stringify(err)}`)
                    this.dispatch({
                        type: "SharingModalErrorAction",
                        error: {
                            type: "adding",
                            message: "Failed to grant access, verify the email address or contact maintainers"
                        }
                    })
                })
        }

        return Promise.resolve()
    }
}