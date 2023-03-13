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

export default function AlbumRouterPage() {
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
