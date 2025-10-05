'use client';

import {createContext, ReactNode, useCallback, useEffect, useMemo, useReducer} from "react";
import {
    AlbumId,
    CatalogFactoryArgs,
    catalogReducer,
    catalogThunks,
    CatalogThunksInterface,
    CatalogViewerAction,
    CatalogViewerState,
    initialCatalogState,
    isRedirectToAlbumIdPayload
} from "../../core/catalog";
import {useApplication, useUnrecoverableErrorDispatch} from "../../core/application";
import {AuthenticatedUser} from "../../core/security";
import {useThunks} from "../../libs/dthunks/react/useThunks";
import {getPayload} from "../../libs/daction";

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    selectedAlbumId?: AlbumId // state managed from the URL
    handlers?: Omit<CatalogThunksInterface, "onPageRefresh">
}

export const CatalogViewerContext = createContext<CatalogViewerStateWithDispatch>({
    state: initialCatalogState({}),
})

export const CatalogViewerProvider = (
    {children, albumId, redirectToAlbumId, authenticatedUser}: {
        albumId?: AlbumId,
        redirectToAlbumId: (albumId: AlbumId) => void
        authenticatedUser: AuthenticatedUser
        children?: ReactNode
    }
) => {
    const unrecoverableErrorDispatch = useUnrecoverableErrorDispatch()

    const [catalog, dispatch] = useReducer(catalogReducer, initialCatalogState(authenticatedUser))
    const dispatchPropagator = useCallback((action: CatalogViewerAction) => {
        dispatch(action)

        const payload = getPayload(action);
        if (isRedirectToAlbumIdPayload(payload) && payload.redirectTo) {
            redirectToAlbumId(payload.redirectTo);
        }
    }, [dispatch, redirectToAlbumId])

    // Use thunks for sharing modal actions instead of ShareController
    const {onPageRefresh, ...thunks} = useCatalogThunks(catalog, dispatchPropagator);

    useEffect(() => {
        onPageRefresh(albumId)
            .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error}));
    }, [onPageRefresh, albumId, unrecoverableErrorDispatch]);

    return (
        <CatalogViewerContext.Provider value={{state: catalog, handlers: thunks, selectedAlbumId: albumId}}>
            {children}
        </CatalogViewerContext.Provider>
    )
}

/**
 * useCatalogThunks aggregates catalog thunks using the generic thunk engine.
 */
function useCatalogThunks(
    state: CatalogViewerState,
    dispatch: (action: CatalogViewerAction) => void
) {
    const app = useApplication();
    const factoryArgs: CatalogFactoryArgs = useMemo(() => ({
        app,
        dispatch,
    }), [app, dispatch]);

    return useThunks(
        catalogThunks,
        factoryArgs,
        state
    );
}

