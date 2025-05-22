import {ShareModal} from "../../catalog";
import {ShareHandlers} from "./CatalogViewerStateWithDispatch";
import {useContext, useMemo} from "react";
import {CatalogViewerContext} from "./CatalogViewerProvider";

/**
 * Hook to access sharing modal state and handlers from CatalogViewerContext.
 */
export function useSharingModalController(): ShareHandlers & {
    shareModal?: ShareModal
} {
    const {state: {shareModal}, handlers} = useContext(CatalogViewerContext);

    // Memoize the handlers to ensure referential stability
    return useMemo(() => ({
        ...handlers,
        shareModal,
    }), [handlers, shareModal]);
}

