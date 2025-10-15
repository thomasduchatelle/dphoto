'use client';

import {Alert, Box, Divider, Drawer, Toolbar} from "@mui/material";
import React from "react";
import AlbumsList from "../AlbumsList";
import MediaList from "../MediasList";
import {AlbumFilterCriterion, AlbumId, CatalogViewerPageSelection} from "../../../core/catalog";
import AlbumListActions from "../AlbumsListActions/AlbumListActions";
import {useClientRouter} from "../../../components/ClientRouter";


export default function MediasPage({
                                       albums,
                                       albumNotFound,
                                       albumsLoaded,
                                       mediasLoaded,
                                       medias,
                                       scrollToMedia,
                                       albumFilterOptions,
                                       albumFilter,
                                       onAlbumFilterChange,
                                       displayedAlbum,
                                       openSharingModal,
                                       openDeleteAlbumDialog,
                                       openEditDatesDialog,
                                       openEditNameDialog,
                                       openCreateDialog,
                                       displayedAlbumIdIsOwned,
                                       deleteButtonEnabled,
                                   }: {
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
    openSharingModal: (albumId: AlbumId) => void
    openDeleteAlbumDialog: () => void
    openEditDatesDialog: () => void
    openCreateDialog: () => void
    openEditNameDialog: () => void
    displayedAlbumIdIsOwned: boolean
    deleteButtonEnabled?: boolean
    scrollToMedia?: string
} & CatalogViewerPageSelection) {
    const drawerWidth = 450
    const {navigate} = useClientRouter();
    
    const handleAlbumClick = (albumId: AlbumId) => {
        navigate(`/albums/${albumId.owner}/${albumId.folderName}`);
    };

    return (
        <Box sx={{display: 'flex'}}>
            <Box
                component="nav"
                sx={{width: {lg: drawerWidth}, flexShrink: {lg: 0}}}
                aria-label="mailbox folders"
            >
                <Drawer
                    variant="permanent"
                    sx={{
                        display: {xs: 'none', lg: 'block'},
                        '& .MuiDrawer-paper': {
                            boxSizing: 'border-box',
                            width: drawerWidth,
                            border: 'none',
                        },
                    }}
                >
                    <Toolbar/>
                    {albumsLoaded && (
                        <>
                            <AlbumListActions
                                selected={albumFilter}
                                options={albumFilterOptions}
                                onAlbumFiltered={onAlbumFilterChange}
                                openDeleteAlbumDialog={openDeleteAlbumDialog}
                                openEditDatesDialog={openEditDatesDialog}
                                openEditNameDialog={openEditNameDialog}
                                openCreateDialog={openCreateDialog}
                                displayedAlbumIdIsOwned={displayedAlbumIdIsOwned}
                                deleteButtonEnabled={deleteButtonEnabled ?? false}
                            />
                            <Divider/>
                        </>
                    )}
                    <AlbumsList albums={albums}
                                loaded={albumsLoaded}
                                selectedAlbumId={displayedAlbum?.albumId}
                                openSharingModal={openSharingModal}
                                onAlbumClick={handleAlbumClick}/>
                </Drawer>
            </Box>
            <Box
                component="main"
                sx={theme => ({
                    padding: theme.spacing(1),
                    flexGrow: 1,
                    width: {lg: `calc(100% - ${drawerWidth}px)`},
                    backgroundColor: theme.palette.background.paper,
                })}
            >
                {(albumsLoaded && !albums && (
                    <Alert severity='info' sx={{mt: 3}}>
                        Your account is empty, start to create new albums and upload your photos with the command line
                        interface.
                    </Alert>
                )) || (
                    <Box>
                        <MediaList
                            medias={medias}
                            loaded={mediasLoaded}
                            albumNotFound={albumNotFound}
                            scrollToMedia={scrollToMedia}
                        />
                    </Box>
                )}
            </Box>
        </Box>
    )
}
