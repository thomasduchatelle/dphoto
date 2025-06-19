import {Box, Toolbar, useMediaQuery, useTheme} from "@mui/material";
import React from 'react';
import AppNav from "../../../components/AppNav";
import UserMenu from "../../../components/user.menu";
import AlbumsList from "./AlbumsList";
import MediasPage from "./MediasPage";
import MobileNavigation from "./MobileNavigation";
import {useAuthenticatedUser, useLogoutCase} from "../../../core/application";
import {useCatalogContext} from "../../../components/catalog-react";
import {useLocation, useSearchParams} from "react-router-dom";
import {albumIdEquals, deleteDialogSelector, sharingDialogSelector} from "../../../core/catalog";
import {CreateAlbumDialogContainer} from "./CreateAlbumDialog";
import AlbumListActions from "./AlbumsListActions/AlbumListActions";
import ShareDialog from "./ShareDialog";
import {DeleteAlbumDialog} from "./DeleteAlbumDialog";

export function CatalogViewerPage() {
    const authenticatedUser = useAuthenticatedUser();

    const {
        state,
        handlers: {
            onAlbumFilterChange,
            createAlbum,
            deleteAlbum,
            closeDeleteAlbumDialog,
            openDeleteAlbumDialog,
            openSharingModal,
            closeSharingModal,
            revokeAlbumAccess,
            grantAlbumSharing,
            openEditAlbumDatesDialog,
        },
        selectedAlbumId
    } = useCatalogContext()
    const logoutCase = useLogoutCase();

    const {pathname} = useLocation()
    const [search] = useSearchParams()

    const theme = useTheme()

    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));
    const isAlbumsPage = pathname === '/albums'

    const selectedAlbum = state.albums.find(album => albumIdEquals(album.albumId, selectedAlbumId))

    return (
        <Box>
            <AppNav
                rightContent={<UserMenu user={authenticatedUser} onLogout={logoutCase.logout}/>}
            />
            <Toolbar/>
            <Box sx={{mt: 2, pl: 2, pr: 2, display: {lg: 'none'}}}>
                <MobileNavigation album={isAlbumsPage ? undefined : selectedAlbum}/>
            </Box>
            <CreateAlbumDialogContainer createAlbum={createAlbum}>
                {(controls) => isMobileDevice && isAlbumsPage ? (
                    <>
                        <AlbumListActions
                            selected={state.albumFilter}
                            options={state.albumFilterOptions}
                            onAlbumFiltered={onAlbumFilterChange}
                            openDeleteAlbumDialog={openDeleteAlbumDialog}
                            openEditAlbumDatesDialog={openEditAlbumDatesDialog}
                            {...controls}
                        />
                        <AlbumsList albums={state.albums}
                                    loaded={state.albumsLoaded}
                                    selectedAlbumId={selectedAlbumId}
                                    openSharingModal={openSharingModal}/>
                    </>
                ) : (
                    <MediasPage
                        {...state}
                        selectedAlbumId={selectedAlbumId}
                        onAlbumFilterChange={onAlbumFilterChange}
                        scrollToMedia={search.get("mediaId") ?? undefined}
                        openSharingModal={openSharingModal}
                        openDeleteAlbumDialog={openDeleteAlbumDialog}
                        openEditAlbumDatesDialog={openEditAlbumDatesDialog}
                        {...controls}
                    />
                )}
            </CreateAlbumDialogContainer>
            <ShareDialog
                {...sharingDialogSelector(state)}
                onClose={closeSharingModal}
                onRevoke={revokeAlbumAccess}
                onGrant={grantAlbumSharing}
            />
            <DeleteAlbumDialog
                {...deleteDialogSelector(state)}
                onDelete={deleteAlbum}
                onClose={closeDeleteAlbumDialog}
            />
        </Box>
    );
}
