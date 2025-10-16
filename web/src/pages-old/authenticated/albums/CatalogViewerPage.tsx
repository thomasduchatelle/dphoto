'use client';

import {Box, Toolbar, useMediaQuery, useTheme} from "@mui/material";
import React from 'react';
import AppNav from "../../../components/AppNav";
import UserMenu from "../../../components/user.menu";
import AlbumsList from "../../../components/albums/AlbumsList";
import MediasPage from "../../../components/albums/MediasPage";
import MobileNavigation from "../../../components/albums/MobileNavigation";
import {useAuthenticatedUser, useLogoutCase} from "../../../core/application";
import {useCatalogContext} from "../../../components/catalog-react";
import {useClientRouter} from "../../../components/ClientRouter";
import {
    AlbumId,
    catalogViewerPageSelector,
    createDialogSelector,
    deleteDialogSelector,
    editDatesDialogSelector,
    editNameDialogSelector,
    sharingDialogSelector,
    albumListActionsSelector
} from "../../../core/catalog";
import {CreateAlbumDialog} from "../../../components/albums/CreateAlbumDialog";
import AlbumListActions from "../../../components/albums/AlbumsListActions/AlbumListActions";
import ShareDialog from "../../../components/albums/ShareDialog";
import {DeleteAlbumDialog} from "../../../components/albums/DeleteAlbumDialog";
import {EditDatesDialog} from "../../../components/albums/EditDatesDialog";
import {EditNameDialog} from "../../../components/albums/EditNameDialog";
import {displayedAlbumSelector} from "../../../core/catalog/language/selector-displayedAlbum";

export function CatalogViewerPage() {
    const authenticatedUser = useAuthenticatedUser();

    const {
        state,
        handlers: {
            onAlbumFilterChange,
            openCreateDialog,
            closeCreateDialog,
            updateCreateDialogStartDate,
            updateCreateDialogEndDate,
            updateCreateDialogStartsAtStartOfTheDay,
            updateCreateDialogEndsAtEndOfTheDay,
            submitCreateAlbum,
            deleteAlbum,
            closeDeleteAlbumDialog,
            openDeleteAlbumDialog,
            openEditDatesDialog,
            closeEditDatesDialog,
            openEditNameDialog,
            closeEditNameDialog,
            openSharingModal,
            closeSharingModal,
            revokeAlbumAccess,
            grantAlbumAccess,
            updateAlbumDates,
            updateEditDatesDialogStartDate,
            updateEditDatesDialogEndDate,
            updateEditDatesDialogStartAtDayStart,
            updateEditDatesDialogEndAtDayEnd,
            saveAlbumName,
            changeFolderNameEnabled,
            changeAlbumName,
            changeFolderName,
        },
        selectedAlbumId
    } = useCatalogContext()
    const logoutCase = useLogoutCase();

    const {path, query, navigate} = useClientRouter()

    const handleAlbumClick = (albumId: AlbumId) => {
        navigate(`/albums/${albumId.owner}/${albumId.folderName}`);
    };

    const theme = useTheme()

    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));
    const isAlbumsPage = path === '/albums'

    const {albumFilter, albumFilterOptions, albumsLoaded, albums, displayedAlbum} = catalogViewerPageSelector(state);

    const editDatesDialogState = editDatesDialogSelector(state);
    const editNameDialogState = editNameDialogSelector(state);
    const createDialogState = createDialogSelector(state);

    return (
        <Box>
            <AppNav
                rightContent={<UserMenu user={authenticatedUser} onLogout={logoutCase.logout}/>}
            />
            <Toolbar/>
            <Box sx={{mt: 2, pl: 2, pr: 2, display: {lg: 'none'}}}>
                <MobileNavigation album={isAlbumsPage ? undefined : displayedAlbum}/>
            </Box>
            {isMobileDevice && isAlbumsPage ? (
                <>
                    <AlbumListActions
                        {...albumListActionsSelector(state)}
                        onAlbumFiltered={onAlbumFilterChange}
                        openDeleteAlbumDialog={openDeleteAlbumDialog}
                        openEditDatesDialog={openEditDatesDialog}
                        openEditNameDialog={openEditNameDialog}
                        openCreateDialog={openCreateDialog}
                    />
                    <AlbumsList albums={albums}
                                loaded={albumsLoaded}
                                selectedAlbumId={selectedAlbumId}
                                openSharingModal={openSharingModal}
                                onAlbumClick={handleAlbumClick}/>
                </>
            ) : (
                <MediasPage
                    {...catalogViewerPageSelector(state)}
                    {...displayedAlbumSelector(state)}
                    deleteButtonEnabled={albumListActionsSelector(state).deleteButtonEnabled}
                    createButtonEnabled={albumListActionsSelector(state).createButtonEnabled}
                    onAlbumFilterChange={onAlbumFilterChange}
                    scrollToMedia={query.get("mediaId") ?? undefined}
                    openSharingModal={openSharingModal}
                    openDeleteAlbumDialog={openDeleteAlbumDialog}
                    openEditDatesDialog={openEditDatesDialog}
                    openEditNameDialog={openEditNameDialog}
                    openCreateDialog={openCreateDialog}
                />
            )}
            <CreateAlbumDialog
                {...createDialogState}
                onClose={closeCreateDialog}
                onSubmit={submitCreateAlbum}
                onNameChange={changeAlbumName}
                onStartDateChange={updateCreateDialogStartDate}
                onEndDateChange={updateCreateDialogEndDate}
                onFolderNameChange={changeFolderName}
                onWithCustomFolderNameChange={changeFolderNameEnabled}
                onStartsAtStartOfTheDayChange={updateCreateDialogStartsAtStartOfTheDay}
                onEndsAtEndOfTheDayChange={updateCreateDialogEndsAtEndOfTheDay}
            />
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
            <EditNameDialog
                {...editNameDialogState}
                onClose={closeEditNameDialog}
                onFolderNameChange={changeFolderName}
                onFolderNameEnabledChange={changeFolderNameEnabled}
                onAlbumNameChange={changeAlbumName}
                onSave={saveAlbumName}
            />
        </Box>
    );
}
