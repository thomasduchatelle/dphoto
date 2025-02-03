import {createContext, ReactNode, useCallback, useEffect, useMemo, useReducer} from "react";
import {
    Album,
    AlbumFilterCriterion,
    AlbumId,
    albumIdEquals,
    catalogReducerFunction,
    CatalogViewerAction,
    initialCatalogState,
    MediaPerDayLoader,
    PostCreateAlbumHandler,
    SelectAlbumHandler
} from "../../catalog";
import {DPhotoApplication, useApplication, useUnrecoverableErrorDispatch} from "../../application";
import {CatalogFactory} from "../../catalog/catalog-factories";
import {CatalogHandlers, CatalogViewerStateWithDispatch} from "./CatalogViewerStateWithDispatch";
import {AuthenticatedUser} from "../../security";
import {AlbumFilterHandler} from "../../catalog/domain/AlbumFilterHandler";

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

        if (!albumId && action.type === "AlbumsAndMediasLoadedAction" && action.selectedAlbum) {
            redirectToAlbumId(action.selectedAlbum.albumId)
        }
        if (action.type === "AlbumsFilteredAction" && action.albumId) {
            redirectToAlbumId(action.albumId)
        }
    }, [dispatch, redirectToAlbumId, albumId])

    const {allAlbums, loadingMediasFor, mediasLoadedFromAlbumId} = catalog
    const handlers = useMemo(
        () => new CompositeHandler(app, dispatchPropagator, allAlbums, albumId),
        [app, dispatchPropagator, allAlbums, albumId]
    )

    useEffect(() => {
        const loader = new CatalogFactory(app).mediaViewLoader()
        loader.loadInitialCatalog({albumId})
            .then(dispatchPropagator)
            .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        const handler = new SelectAlbumHandler(dispatchPropagator, new MediaPerDayLoader(new CatalogFactory(app).restAdapter()), mediasLoadedFromAlbumId, loadingMediasFor)
        if (albumId) {
            handler.onSelectAlbum(albumId)
                .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))

        }
    }, [app, dispatchPropagator, unrecoverableErrorDispatch, albumId, mediasLoadedFromAlbumId, loadingMediasFor]);

    return (
        <CatalogViewerContext.Provider value={{state: catalog, handlers, selectedAlbumId: albumId}}>
            {children}
        </CatalogViewerContext.Provider>
    )
}

class CompositeHandler implements CatalogHandlers {
    constructor(
        private readonly app: DPhotoApplication,
        private readonly dispatch: (action: CatalogViewerAction) => void,
        private readonly allAlbums: Album[],
        private readonly albumId?: AlbumId,
    ) {
        const selectedAlbum = allAlbums.find(album => albumIdEquals(albumId, album.albumId))

        this.onAlbumFilterChange = new AlbumFilterHandler(dispatch, {selectedAlbum, allAlbums}).onAlbumFilter
        this.onAlbumCreated = new PostCreateAlbumHandler(dispatch, new CatalogFactory(app).mediaViewLoader()).onAlbumCreated
    }

    onAlbumCreated: (albumId: AlbumId) => Promise<void>
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}

