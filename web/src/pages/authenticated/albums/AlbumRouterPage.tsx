import {Box, Toolbar, useMediaQuery, useTheme} from "@mui/material";
import React, {useCallback} from 'react';
import AppNav from "../../../components/AppNav";
import UserMenu from "../../../components/user.menu";
import AlbumsList from "./AlbumsList";
import MediasPage from "./MediasPage";
import MobileNavigation from "./MobileNavigation";
import {useAuthenticatedUser, useLogoutCase} from "../../../core/application";
import {AlbumId, useCatalogController} from "../../../core/catalog";
import {useLocation} from "react-router-dom";

// type AlbumPageUrlParams = {
//     owner: string | undefined,
//     album: string | undefined,
// }

export default function AlbumRouterPage() {
    // const mustBeAuthenticated = useMustBeAuthenticated()
    // const {owner, album} = useParams<AlbumPageUrlParams>()
    // const [state, setState] = useState<CatalogStateV0>({
    //   fullyLoaded: false,
    //   albumsLoaded: false,
    //   albumNotFound: false,
    //   albums: [],
    //   medias: [],
    // })
    // const cache = useRef<AlbumsLogicCache>({owner: '', albums: []})
    // const navigate = useNavigate()
    // const {width} = useWindowDimensions()
    //
    // const engine = useMemo(() => new AlbumsLogic(mustBeAuthenticated, cache, width,
    //   {
    //     redirectToAlbum(albumId: AlbumId): void {
    //       navigate(`/albums/${albumId.owner}/${albumId.folderName}`, {replace: false})
    //     },
    //     renderAlbumNotPresent(albums: Album[], albumId: AlbumId): void {
    //       setState({
    //         albums,
    //         fullyLoaded: true,
    //         albumsLoaded: true,
    //         albumNotFound: true,
    //         medias: [],
    //         selectedAlbum: {
    //           albumId: albumId,
    //           end: new Date(),
    //           name: "not found",
    //           relativeTemperature: 0,
    //           start: new Date(),
    //           temperature: 0,
    //           totalCount: 0
    //         },
    //       })
    //     },
    //     renderAlbumsList(albums: Album[]): void {
    //       setState(current => ({...current, albums, albumsLoaded: true, fullyLoaded: false}))
    //     },
    //     renderAlbumsWithMedia(albums: Album[], selectedAlbum: Album, medias: MediaWithinADay[]): void {
    //       setState(current => ({
    //         ...current,
    //         albums,
    //         medias,
    //         selectedAlbum,
    //         fullyLoaded: true,
    //         albumsLoaded: true,
    //         albumNotFound: false,
    //       }))
    //     },
    //     renderNoAlbums(): void {
    //       setState(current => ({
    //         ...current,
    //         fullyLoaded: true,
    //         albumsLoaded: true,
    //         albums: [],
    //         medias: [],
    //         albumNotFound: false,
    //       }))
    //     }
    //   }), [mustBeAuthenticated, cache, width, setState, navigate]) //mustBeAuthenticated
    //
    // useEffect(() => {
    //   if (!owner || !album) {
    //     engine.loadAlbumsPage().catch(err => console.log(`Error: ${err}\n${err.stack}`))
    //   } else {
    //     engine.loadMediasPage({
    //       owner,
    //       folderName: album
    //     }).catch(err => console.log(`Error: ${err}\n${err.stack}`))
    //   }
    // }, [engine, owner, album])
    const {albums, selectedAlbum, albumNotFound, medias} = useCatalogController()
    const authenticatedUser = useAuthenticatedUser();
    const logoutCase = useLogoutCase();

    const selectAlbum = useCallback((selected: AlbumId) => console.log(`Selected: ${selected}`), [])

    const {pathname} = useLocation()
    const theme = useTheme()

    // '/albums' page is only available on small devices
    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));
    const isAlbumsPage = pathname === '/albums'

    return (
        <Box>
            <AppNav
                rightContent={<UserMenu user={authenticatedUser} onLogout={logoutCase.logout}/>}
            />
            <Toolbar/>
            <Box sx={{mt: 2, pl: 2, pr: 2, display: {lg: 'none'}}}>
                <MobileNavigation album={selectedAlbum}/>
            </Box>
            {isMobileDevice && isAlbumsPage ? (
                <AlbumsList albums={albums}
                            loaded={true}
                            selected={selectedAlbum}/>
            ) : (
                <MediasPage
                    albums={albums}
                    albumNotFound={albumNotFound}
                    fullyLoaded={true}
                    medias={medias}
                    selectAlbum={selectAlbum}
                    selectedAlbum={selectedAlbum}
                />
            )}
        </Box>
    );
}
