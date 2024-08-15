import {createContext, ReactNode, useEffect, useReducer} from "react";
import {catalogReducer, initialCatalogState} from "./catalog-reducer";
import {useApplication, useUnrecoverableErrorDispatch} from "../../application";
import {CatalogFactory} from "../../catalog/catalog-factories";
import {CatalogViewerAction, CatalogViewerStateWithDispatch} from "./catalog-viewer-state";
import {AlbumId} from "../../catalog";

export const CatalogViewerContext = createContext<CatalogViewerStateWithDispatch>({
    state: initialCatalogState, dispatch: () => {
    }
})

export const CatalogViewerProvider = (
    {children, albumId, onSelectedAlbumIdByDefault}: {
        albumId?: AlbumId,
        onSelectedAlbumIdByDefault: (albumId: AlbumId) => void
        children?: ReactNode
    }
) => {
    const app = useApplication()
    const unrecoverableErrorDispatch = useUnrecoverableErrorDispatch()

    const [catalog, dispatch] = useReducer(catalogReducer, initialCatalogState)
    const {albumsLoaded, selectedAlbum} = catalog

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
        if (albumId) {
            const selectAlbumHandler = new CatalogFactory(app).selectAlbumHandler()
            selectAlbumHandler.onSelectAlbum({
                albumId: albumId,
                currentAlbumId: selectedAlbum?.albumId,
                loaded: albumsLoaded,
            }, (action: CatalogViewerAction) => dispatch(action))
                .catch(error => unrecoverableErrorDispatch({type: 'unrecoverable-error', error: error}))

        }
    }, [albumId, albumsLoaded, selectedAlbum, app, unrecoverableErrorDispatch]);

    return (
        <CatalogViewerContext.Provider value={{state: catalog, dispatch}}>
            {children}
        </CatalogViewerContext.Provider>
    )
}
