import {Box, Toolbar, useMediaQuery, useTheme} from "@mui/material";
import React from 'react';
import AppNav from "../../../components/AppNav";
import UserMenu from "../../../components/user.menu";
import AlbumsList from "./AlbumsList";
import MediasPage from "./MediasPage";
import MobileNavigation from "./MobileNavigation";
import {useAuthenticatedUser, useLogoutCase} from "../../../core/application";
import {useCatalogViewerState} from "../../../core/catalog-react";
import {useLocation, useSearchParams} from "react-router-dom";

export function CatalogViewerPage() {
    const authenticatedUser = useAuthenticatedUser();

    const {
        albums,
        selectedAlbum,
        albumNotFound,
        medias,
        albumsLoaded,
        mediasLoaded
    } = useCatalogViewerState()
    const logoutCase = useLogoutCase();

    const {pathname} = useLocation()
    const [search] = useSearchParams()

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
                <MobileNavigation album={isAlbumsPage ? undefined : selectedAlbum}/>
            </Box>
            {isMobileDevice && isAlbumsPage ? (
                <AlbumsList albums={albums}
                            loaded={albumsLoaded}
                            selected={selectedAlbum}/>
            ) : (
                <MediasPage
                    albums={albums}
                    albumNotFound={albumNotFound}
                    albumsLoaded={albumsLoaded}
                    mediasLoaded={mediasLoaded}
                    medias={medias}
                    selectedAlbum={selectedAlbum}
                    scrollToMedia={search.get("mediaId") ?? undefined}
                />
            )}
        </Box>
    );
}
