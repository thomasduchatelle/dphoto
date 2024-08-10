import {useMatch, useNavigate, useParams} from "react-router-dom";
import {CatalogController, InitialCatalogController} from "./catalog-controller";
import {useCallback, useContext, useEffect, useMemo} from "react";
import {AuthenticatedUser} from "../security";
import {useApplication, useUnrecoverableErrorDispatch} from "../application";
import {CatalogAPIAdapter} from "../catalog-apis";
import {CatalogContext} from "./catalog-context";
import {CatalogState} from "../catalog/catalog-model";


export const useCatalogLoader = (): (user: AuthenticatedUser) => Promise<void> => {
    const match = useMatch('/albums/:owner/:folderName');
    const navigate = useNavigate()
    const app = useApplication()
    const {dispatch} = useContext(CatalogContext)
    const unrecoverableErrorDispatch = useUnrecoverableErrorDispatch()

    const controller = useMemo(() => {
        console.log("new InitialCatalogController")
        return new InitialCatalogController(new CatalogAPIAdapter(app.axiosInstance, app), dispatch, unrecoverableErrorDispatch)
    }, [app, dispatch, unrecoverableErrorDispatch])

    return useCallback(user => {
        return controller.loadInitialCatalog(user.email, match?.params.owner, match?.params.folderName)
            .then(redirect => {
                if (redirect.albumPage) {
                    navigate('/albums')
                } else if (redirect.albumId) {
                    navigate(`/albums/${redirect.albumId.owner}/${redirect.albumId.folderName}`)
                }
            })
            .catch(() => navigate('/albums')) // error will be displayed on main page...
    }, [match, controller, navigate])
}

type AlbumPageUrlParams = {
    owner: string | undefined,
    album: string | undefined,
}

export type CatalogContextWithController = CatalogState

export const useCatalogController = (): CatalogContextWithController => {
    const {catalog, dispatch} = useContext(CatalogContext)
    const {owner, album} = useParams<AlbumPageUrlParams>()
    const app = useApplication()
    const controller = useMemo(() => new CatalogController(new CatalogAPIAdapter(app.axiosInstance, app), dispatch), [app, dispatch])

    useEffect(() => {
        if (owner && album && (owner !== catalog.selectedAlbum?.albumId.owner || album !== catalog.selectedAlbum?.albumId.folderName)) {
            controller.selectAlbum({owner: owner, folderName: album}).then()
        }
    }, [owner, album, catalog.selectedAlbum, controller])

    return catalog
}