import {Alert, Box, Drawer, Toolbar} from "@mui/material";
import React from "react";
import AlbumsList from "../AlbumsList";
import MediaList from "../MediasList";
import {Album, MediaWithinADay} from "../../../../core/catalog";

export default function MediasPage({
                                       albums,
                                       albumNotFound,
                                       fullyLoaded,
                                       medias,
                                       selectedAlbum,
                                   }: {
    albums: Album[]
    albumNotFound: boolean
    fullyLoaded: boolean
    medias: MediaWithinADay[]
    selectedAlbum?: Album
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
                    <AlbumsList albums={albums} loaded={fullyLoaded} selected={selectedAlbum}/>
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
                {(fullyLoaded && !albums && (
                    <Alert severity='info' sx={{mt: 3}}>
                        Your account is empty, start to create new albums and upload your photos with the command line
                        interface.
                    </Alert>
                )) || (
                    <Box>
                        <MediaList medias={medias} loaded={fullyLoaded}
                                   albumNotFound={albumNotFound}/>
                    </Box>
                )}
            </Box>
        </Box>
    )
}