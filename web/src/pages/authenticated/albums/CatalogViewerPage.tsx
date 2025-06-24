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
import {catalogViewerPageSelector, deleteDialogSelector, editDatesDialogSelector, sharingDialogSelector} from "../../../core/catalog";
import {CreateAlbumDialogContainer} from "./CreateAlbumDialog";
import AlbumListActions from "./AlbumsListActions/AlbumListActions";
import ShareDialog from "./ShareDialog";
import {DeleteAlbumDialog} from "./DeleteAlbumDialog";
import {EditDatesDialog} from "./EditDatesDialog";
import {displayedAlbumSelector} from "../../../core/catalog/language/selector-displayedAlbum";

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
            openEditDatesDialog,
            closeEditDatesDialog,
            openSharingModal,
            closeSharingModal,
            revokeAlbumAccess,
            grantAlbumAccess,
            updateAlbumDates,
            updateEditDatesDialogStartDate,
            updateEditDatesDialogEndDate,
            updateEditDatesDialogStartAtDayStart,
            updateEditDatesDialogEndAtDayEnd,
        },
        selectedAlbumId
    } = useCatalogContext()
    const logoutCase = useLogoutCase();

    const {pathname} = useLocation()
    const [search] = useSearchParams()

    const theme = useTheme()

    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));
    const isAlbumsPage = pathname === '/albums'

    const {albumFilter, albumFilterOptions, albumsLoaded, albums, displayedAlbum} = catalogViewerPageSelector(state);

    const editDatesDialogState = editDatesDialogSelector(state);

    return (
        <Box>
            <AppNav
                rightContent={<UserMenu user={authenticatedUser} onLogout={logoutCase.logout}/>}
            />
            <Toolbar/>
            <Box sx={{mt: 2, pl: 2, pr: 2, display: {lg: 'none'}}}>
                <MobileNavigation album={isAlbumsPage ? undefined : displayedAlbum}/>
            </Box>
            <CreateAlbumDialogContainer createAlbum={createAlbum}>
                {(controls) => isMobileDevice && isAlbumsPage ? (
                    <>
                        <AlbumListActions
                            selected={albumFilter}
                            options={albumFilterOptions}
                            onAlbumFiltered={onAlbumFilterChange}
                            openDeleteAlbumDialog={openDeleteAlbumDialog}
                            openEditDatesDialog={openEditDatesDialog}
                            {...displayedAlbumSelector(state)}
                            {...controls}
                        />
                        <AlbumsList albums={albums}
                                    loaded={albumsLoaded}
                                    selectedAlbumId={selectedAlbumId}
                                    openSharingModal={openSharingModal}/>
                    </>
                ) : (
                    <MediasPage
                        {...catalogViewerPageSelector(state)}
                        {...displayedAlbumSelector(state)}
                        onAlbumFilterChange={onAlbumFilterChange}
                        scrollToMedia={search.get("mediaId") ?? undefined}
                        openSharingModal={openSharingModal}
                        openDeleteAlbumDialog={openDeleteAlbumDialog}
                        openEditDatesDialog={openEditDatesDialog}
                        {...controls}
                    />
                )}
            </CreateAlbumDialogContainer>
            <ShareDialog
                {...sharingDialogSelector(state)}
                onClose={closeSharingModal}
                onRevoke={revokeAlbumAccess}
                onGrant={grantAlbumAccess}
            />
            <DeleteAlbumDialog
                {...deleteDialogSelector(state)}
                onDelete={deleteAlbum}
                onClose={closeDeleteAlbumDialog}
            />
            <EditDatesDialog
                {...editDatesDialogState}
                onClose={closeEditDatesDialog}
                onSave={updateAlbumDates}
                onStartDateChange={updateEditDatesDialogStartDate}
                onEndDateChange={updateEditDatesDialogEndDate}
                onStartAtDayStartChange={updateEditDatesDialogStartAtDayStart}
                onEndAtDayEndChange={updateEditDatesDialogEndAtDayEnd}
            />
        </Box>
    );
}
