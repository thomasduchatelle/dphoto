import {createContext, ReactNode, useCallback, useEffect, useReducer} from "react";
import {catalogReducerFunction, initialCatalogState} from "../../catalog/domain/catalog-reducer";
import {useApplication, useUnrecoverableErrorDispatch} from "../../application";
import {CatalogFactory} from "../../catalog/catalog-factories";
import {AlbumFilterCriterion, AlbumId, albumIdEquals, MediaPerDayLoader, SelectAlbumHandler} from "../../catalog";
import {CatalogViewerStateWithDispatch} from "./CatalogViewerStateWithDispatch";
import {AuthenticatedUser} from "../../security";
import {AlbumFilterHandler, AlbumFilterHandlerDispatch} from "../../catalog/domain/AlbumFilterHandler";

export const CatalogViewerContext = createContext<CatalogViewerStateWithDispatch>({
    state: initialCatalogState,
    dispatch: () => {
    },
    handlers: {
        onAlbumFilterChange: () => {
        }
    }
})

export const CatalogViewerProvider = (
    {children, albumId, onSelectedAlbumIdByDefault, authenticatedUser}: {
        albumId?: AlbumId,
        onSelectedAlbumIdByDefault: (albumId: AlbumId) => void
        authenticatedUser: AuthenticatedUser
        children?: ReactNode
    }
) => {
    const app = useApplication()
    const unrecoverableErrorDispatch = useUnrecoverableErrorDispatch()

    const [catalog, dispatch] = useReducer(catalogReducerFunction(authenticatedUser), initialCatalogState)
    const {albumsLoaded, allAlbums} = catalog

    const onAlbumFilterChange = useCallback((criterion: AlbumFilterCriterion) => {
        const selectedAlbum = allAlbums.find(album => albumIdEquals(albumId, album.albumId))
        const overriddenDispatch: AlbumFilterHandlerDispatch = (action) => {
            if (action.type === "AlbumsFilteredAction" && action.albumId) {
                onSelectedAlbumIdByDefault(action.albumId)
            }
            dispatch(action)
        }

        return new AlbumFilterHandler(overriddenDispatch, {selectedAlbum, allAlbums}).onAlbumFilter(criterion)

    }, [dispatch, onSelectedAlbumIdByDefault, albumId, allAlbums])

    const handlers: CatalogHandlers = {
        onAlbumFilterChange,
    }

    useEffect(() => {
        const loader = new CatalogFactory(app).mediaViewLoader()
        loader.loadInitialCatalog({albumId})
            .then(action => {
                dispatch(action);

                if (!albumId && action.type === "AlbumsAndMediasLoadedAction" && action.selectedAlbum) {
                    onSelectedAlbumIdByDefault(action.selectedAlbum.albumId)
                }
            })
            .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        const handler = new SelectAlbumHandler(new MediaPerDayLoader(new CatalogFactory(app).restAdapter()));
        if (albumId) {
            handler.onSelectAlbum({
                albumId: albumId,
                currentAlbumId: undefined,
                loaded: albumsLoaded,
            }, dispatch)
                .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))

        }
    }, [albumId, albumsLoaded, app, unrecoverableErrorDispatch]);

    return (
        <CatalogViewerContext.Provider value={{state: catalog, dispatch, handlers, selectedAlbumId: albumId}}>
            {children}
        </CatalogViewerContext.Provider>
    )
}

export interface CatalogHandlers {
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}