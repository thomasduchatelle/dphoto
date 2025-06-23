export * from "./action-albumAccessGranted"
export * from "./action-sharingModalClosed"
export * from "./action-sharingModalErrorOccurred"
export * from "./action-sharingModalOpened"
export * from "./action-albumAccessRevoked"
export * from "./selector-sharingDialogSelector"
import {closeSharingModalDeclaration} from "./thunk-closeSharingModal"
import {grantAlbumAccessDeclaration} from "./thunk-grantAlbumAccess"
import {openSharingModalDeclaration} from "./thunk-openSharingModal"
import {revokeAlbumAccessDeclaration} from "./thunk-revokeAlbumAccess"

/**
 * Thunks related to album sharing.
 *
 * Expected handler types:
 * - `openSharingModal`: `(album: Album) => void`
 * - `closeSharingModal`: `() => void`
 * - `revokeAlbumAccess`: `(email: string) => Promise<void>`
 * - `grantAlbumAccess`: `(email: string) => Promise<void>`
 */
export const sharingThunks = {
    openSharingModal: openSharingModalDeclaration,
    closeSharingModal: closeSharingModalDeclaration,
    revokeAlbumAccess: revokeAlbumAccessDeclaration,
    grantAlbumAccess: grantAlbumAccessDeclaration,
}
