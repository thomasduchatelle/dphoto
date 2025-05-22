import {createContext, ReactNode, useCallback, useEffect, useMemo, useReducer} from "react";
import {
    Album,
    AlbumFilterCriterion,
    AlbumFilterHandler,
    AlbumId,
    albumIdEquals,
    CatalogLoader,
    catalogReducerFunction,
    CatalogViewerAction,
    initialCatalogState,
    isRedirectToAlbumIdAction,
    MediaPerDayLoader,
    PostCreateAlbumHandler,
    SharingType
} from "../../catalog";
import {DPhotoApplication, useApplication, useUnrecoverableErrorDispatch} from "../../application";
import {CatalogFactory} from "../../catalog/catalog-factories";
import {CatalogHandlers, CatalogViewerStateWithDispatch, ShareHandlers} from "./CatalogViewerStateWithDispatch";
import {AuthenticatedUser} from "../../security";
import {ShareController} from "../../catalog/domain/ShareController";
import {CatalogAPIAdapter} from "../../catalog/adapters/api";

export const CatalogViewerContext = createContext<CatalogViewerStateWithDispatch>({
    state: initialCatalogState,
    handlers: {
        onAlbumFilterChange: () => {
        },
        async onAlbumCreated(albumId: AlbumId): Promise<void> {
        },
        async onRevoke(email: string): Promise<void> {
        },
        async onGrant(email: string, role: SharingType): Promise<void> {
        },
        openSharingModal(album: Album): void {
        },
        onClose(): void {
        },
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

    const {allAlbums, mediasLoadedFromAlbumId, albumsLoaded, loadingMediasFor, shareModal} = catalog

    const handlers = useMemo(() => {
            const sharingAPI = new CatalogAPIAdapter(app.axiosInstance, app);
            const shareController = new ShareController(dispatchPropagator, sharingAPI);

            return new CompositeHandler(app, dispatchPropagator, shareController, allAlbums, albumId, shareModal?.sharedAlbumId)
        },
        [app, dispatchPropagator, allAlbums, albumId, shareModal]
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

class CompositeHandler implements CatalogHandlers, ShareHandlers {
    constructor(
        app: DPhotoApplication,
        dispatch: (action: CatalogViewerAction) => void,
        private readonly shareController: ShareController,
        allAlbums: Album[],
        albumId?: AlbumId,
        private readonly sharingModalAlbumId?: AlbumId,
    ) {
        const selectedAlbum = allAlbums.find(album => albumIdEquals(albumId, album.albumId))

        this.onAlbumFilterChange = new AlbumFilterHandler(dispatch, {selectedAlbum, allAlbums}).onAlbumFilter
        this.onAlbumCreated = new PostCreateAlbumHandler(dispatch, new CatalogFactory(app).restAdapter()).onAlbumCreated
    }

    onAlbumCreated: (albumId: AlbumId) => Promise<void>
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void

    onRevoke = async (email: string): Promise<void> => {
        if (!this.sharingModalAlbumId) {
            return Promise.reject("No album selected");
        }
        return this.shareController.revokeAccess(this.sharingModalAlbumId, email);
    }

    onGrant = async (email: string, role: SharingType): Promise<void> => {
        if (!this.sharingModalAlbumId) {
            return Promise.reject("No album selected");
        }
        return this.shareController.grantAccess(this.sharingModalAlbumId, email, role);
    }

    openSharingModal = (album: Album): void => {
        this.shareController.openSharingModal(album);
    }

    onClose = (): void => {
        this.shareController.onClose();
    }
}

