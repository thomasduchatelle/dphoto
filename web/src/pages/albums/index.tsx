import {Box, Toolbar, useMediaQuery, useTheme} from "@mui/material";
import React from 'react';
import AppNav from "../AppNav";
import UserMenu from "../user.menu";
import AlbumsList from "./AlbumsList";
import MediasPage from "./MediasPage";
import MobileNavigation from "./MobileNavigation";
import {useAuthenticatedUser, useLogoutCase} from "../../core/application";
import {useCatalogContext} from "../catalog-react";
import {useLocation, useSearchParams} from "react-router-dom";
import {
    catalogViewerPageSelector,
    createDialogSelector,
    deleteDialogSelector,
    editDatesDialogSelector,
    editNameDialogSelector,
    sharingDialogSelector
} from "../../core/catalog";
import {CreateAlbumDialog} from "./CreateAlbumDialog";
import AlbumListActions from "./AlbumsListActions/AlbumListActions";
import ShareDialog from "./ShareDialog";
import {DeleteAlbumDialog} from "./DeleteAlbumDialog";
import {EditDatesDialog} from "./EditDatesDialog";
import {EditNameDialog} from "./EditNameDialog";
import {displayedAlbumSelector} from "../../core/catalog/language/selector-displayedAlbum";

export default function AlbumsIndexPage() {
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

    const {pathname} = useLocation()
    const [search] = useSearchParams()

    const theme = useTheme()

    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));
    const isAlbumsPage = pathname === '/albums'

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
                        selected={albumFilter}
                        options={albumFilterOptions}
                        onAlbumFiltered={onAlbumFilterChange}
                        openDeleteAlbumDialog={openDeleteAlbumDialog}
                        openEditDatesDialog={openEditDatesDialog}
                        openEditNameDialog={openEditNameDialog}
                        openCreateDialog={openCreateDialog}
                        {...displayedAlbumSelector(state)}
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
