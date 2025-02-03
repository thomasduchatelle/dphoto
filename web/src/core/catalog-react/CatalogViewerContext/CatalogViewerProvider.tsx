import {createContext, ReactNode, useCallback, useEffect, useMemo, useReducer} from "react";
import {
    Album,
    AlbumFilterCriterion,
    AlbumId,
    albumIdEquals,
    catalogReducerFunction,
    CatalogViewerAction,
    initialCatalogState,
    isRedirectToAlbumIdAction,
    MediaPerDayLoader,
    PostCreateAlbumHandler
} from "../../catalog";
import {DPhotoApplication, useApplication, useUnrecoverableErrorDispatch} from "../../application";
import {CatalogFactory} from "../../catalog/catalog-factories";
import {CatalogHandlers, CatalogViewerStateWithDispatch} from "./CatalogViewerStateWithDispatch";
import {AuthenticatedUser} from "../../security";
import {AlbumFilterHandler} from "../../catalog/domain/AlbumFilterHandler";
import {CatalogLoader} from "../../catalog/domain/CatalogLoader";

export const CatalogViewerContext = createContext<CatalogViewerStateWithDispatch>({
    state: initialCatalogState,
    handlers: {
        onAlbumFilterChange: () => {
        },
        async onAlbumCreated(albumId: AlbumId): Promise<void> {
        }
    }
})

export const CatalogViewerProvider = (
    {children, albumId, redirectToAlbumId, authenticatedUser}: {
        albumId?: AlbumId,
        redirectToAlbumId: (albumId: AlbumId) => void
        authenticatedUser: AuthenticatedUser
        children?: ReactNode
    }
) => {
    const app = useApplication()
    const unrecoverableErrorDispatch = useUnrecoverableErrorDispatch()

    const [catalog, dispatch] = useReducer(catalogReducerFunction(authenticatedUser), initialCatalogState)
    const dispatchPropagator = useCallback((action: CatalogViewerAction) => {
        dispatch(action)

        if (isRedirectToAlbumIdAction(action) && action.redirectTo) {
            redirectToAlbumId(action.redirectTo);
        }
    }, [dispatch, redirectToAlbumId])

    const {allAlbums, mediasLoadedFromAlbumId, albumsLoaded, loadingMediasFor} = catalog
    const handlers = useMemo(
        () => new CompositeHandler(app, dispatchPropagator, allAlbums, albumId),
        [app, dispatchPropagator, allAlbums, albumId]
    )

    useEffect(() => {
        const restAdapter = new CatalogFactory(app).restAdapter();
        const loader = new CatalogLoader(dispatchPropagator, new MediaPerDayLoader(restAdapter), restAdapter, {
            mediasLoadedFromAlbumId,
            allAlbums,
            albumsLoaded,
            loadingMediasFor,
        })
        loader.onPageRefresh(albumId)
            .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))
    }, [app, dispatchPropagator, mediasLoadedFromAlbumId, allAlbums, albumsLoaded, albumId, loadingMediasFor, unrecoverableErrorDispatch])

    return (
        <CatalogViewerContext.Provider value={{state: catalog, handlers, selectedAlbumId: albumId}}>
            {children}
        </CatalogViewerContext.Provider>
    )
}

class CompositeHandler implements CatalogHandlers {
    constructor(
        app: DPhotoApplication,
        dispatch: (action: CatalogViewerAction) => void,
        allAlbums: Album[],
        albumId?: AlbumId,
    ) {
        const selectedAlbum = allAlbums.find(album => albumIdEquals(albumId, album.albumId))

        this.onAlbumFilterChange = new AlbumFilterHandler(dispatch, {selectedAlbum, allAlbums}).onAlbumFilter
        this.onAlbumCreated = new PostCreateAlbumHandler(dispatch, new CatalogFactory(app).restAdapter()).onAlbumCreated
    }

    onAlbumCreated: (albumId: AlbumId) => Promise<void>
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}

