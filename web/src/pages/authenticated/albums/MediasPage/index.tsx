import {Alert, Box, Divider, Drawer, Toolbar} from "@mui/material";
import React from "react";
import AlbumsList from "../AlbumsList";
import MediaList from "../MediasList";
import {Album, AlbumFilterCriterion, AlbumFilterEntry, AlbumId, CreateAlbumControls, MediaWithinADay} from "../../../../core/catalog";
import AlbumListActions from "../AlbumsListActions/AlbumListActions";
import {CatalogViewerPageSelection} from "../../../../core/catalog/navigation/selector-catalog-viewer-page";

const albumFilterFeatureFlag = true

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
                                       selectedAlbumId,
                                       openSharingModal,
                                       openDeleteAlbumDialog,
                                       openEditDatesDialog,
                                       ...controls
                                   }: {
    selectedAlbumId: AlbumId | undefined
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
    openSharingModal: (album: Album) => void
    openDeleteAlbumDialog: () => void
    openEditDatesDialog: () => void
} & CreateAlbumControls & CatalogViewerPageSelection) {
    const drawerWidth = 450

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
                    {albumFilterFeatureFlag && albumsLoaded && (
                        <>
                            <AlbumListActions
                                selected={albumFilter}
                                options={albumFilterOptions}
                                onAlbumFiltered={onAlbumFilterChange}
                                openDeleteAlbumDialog={openDeleteAlbumDialog}
                                openEditDatesDialog={openEditDatesDialog}
                                {...controls}
                            />
                            <Divider/>
                        </>
                    )}
                    <AlbumsList albums={albums}
                                loaded={albumsLoaded}
                                selectedAlbumId={selectedAlbumId}
                                openSharingModal={openSharingModal}/>
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
