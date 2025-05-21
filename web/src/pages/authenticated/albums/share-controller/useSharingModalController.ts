import {useMemo, useReducer} from "react";
import {ShareController} from "./ShareController";
import {ShareState, sharingModalReducer} from "./sharingModalReducer";
import {useCatalogAPIAdapter} from "../../../../core/catalog-react";
import {Album, SharingType} from "../../../../core/catalog";

export interface ShareCallbacks {

    onRevoke(email: string): Promise<void>

    onGrant(email: string, role: SharingType): Promise<void>

    openSharingModal(album: Album): void

    onClose(): void
}

export function useSharingModalController(): ShareState & ShareCallbacks {
    const [state, dispatch] = useReducer(sharingModalReducer, {
        open: false,
        sharedWith: [],
    })
    const catalogAPIAdapter = useCatalogAPIAdapter()

    const callbacks = useMemo(() => {
        const ctrl = new ShareController(dispatch, catalogAPIAdapter);
        return {
            ...ctrl,
            onGrant: (email: string, role: SharingType) => {
                if (!state.sharedAlbumId) {
                    return Promise.reject(`ERROR: no albumId selected to be granted, cannot grant access for ${email}`)
                }

                return ctrl.grantAccess(state.sharedAlbumId, email, role)
            },
            onRevoke: (email: string) => {
                if (!state.sharedAlbumId) {
                    return Promise.reject(`ERROR: no albumId selected to be revoked, cannot revoke access for ${email}`)
                }

                return ctrl.revokeAccess(state.sharedAlbumId, email)
            },
        }
    }, [dispatch, catalogAPIAdapter, state.sharedAlbumId])

    return {
        ...callbacks,
        ...state
    }
}

