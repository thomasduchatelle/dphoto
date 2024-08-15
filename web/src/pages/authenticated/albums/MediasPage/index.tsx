import {Alert, Box, Divider, Drawer, Toolbar} from "@mui/material";
import React from "react";
import AlbumsList from "../AlbumsList";
import MediaList from "../MediasList";
import AlbumListActions from "../AlbumsListActions";
import {Album, MediaWithinADay} from "../../../../core/catalog";

const albumFilterFeature = false

export default function MediasPage({
                                       albums,
                                       albumNotFound,
                                       albumsLoaded,
                                       mediasLoaded,
                                       medias,
                                       selectedAlbum,
                                       scrollToMedia,
                                   }: {
    albums: Album[]
    albumNotFound: boolean
    albumsLoaded: boolean
    mediasLoaded: boolean
    medias: MediaWithinADay[]
    selectedAlbum?: Album
    scrollToMedia?: string
}) {
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
                    {albumFilterFeature && albumsLoaded && (
                        <>
                            <AlbumListActions
                                selected={{
                                    criterion: {owners: []},
                                    avatars: [],
                                    name: "All Albums",
                                }}
                                options={[]}
                                onAlbumFiltered={() => {
                                }}
                            />
                            <Divider/>
                        </>
                    )}
                    <AlbumsList albums={albums} loaded={albumsLoaded} selected={selectedAlbum}/>
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