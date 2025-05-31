import {Album, ShareModal, SharingType} from "../../catalog";
import {useMemo} from "react";
import {useCatalogContext} from "./useCatalogContext";

export interface ShareHandlers {

    onRevoke(email: string): Promise<void>

    onGrant(email: string, role: SharingType): Promise<void>

    openSharingModal(album: Album): void

    onClose(): void
}

/**
 * Hook to access sharing modal state and handlers from CatalogViewerContext.
 */
export function useSharingModalController(): ShareHandlers & {
    shareModal?: ShareModal
} {
    const {state: {shareModal}, handlers} = useCatalogContext();

    // Memoize the handlers to ensure referential stability
    return useMemo(() => ({
        onClose: handlers.closeSharingModal,
        onRevoke: handlers.revokeAlbumSharing,
        onGrant: handlers.grantAlbumSharing,
        openSharingModal: handlers.openSharingModal,
        shareModal,
    }), [handlers, shareModal]);
}

