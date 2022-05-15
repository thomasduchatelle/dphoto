import {Alert, Box, Drawer, Toolbar} from "@mui/material";
import React from "react";
import AlbumsList from "../AlbumsList";
import {Album, AlbumId, MediaWithinADay} from "../logic";
import MediaList from "../MediasList";

export default function MediasPage({
                                     albums,
                                     albumNotFound,
                                     fullyLoaded,
                                     medias,
                                     selectAlbum,
                                     selectedAlbum,
                                   }: {
  albums: Album[]
  albumNotFound: boolean
  fullyLoaded: boolean
  medias: MediaWithinADay[]
  selectAlbum(albumId: AlbumId): void
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
            width: drawerWidth,
            display: {xs: 'none', lg: 'flex'},
            flexShrink: 0,
            [`& .MuiDrawer-paper`]: {width: drawerWidth, boxSizing: 'border-box'},
          }}
        >
          <Toolbar/>
          <AlbumsList albums={albums} loaded={fullyLoaded} selected={selectedAlbum}
                      onSelect={selectAlbum}/>
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{flexGrow: 1, pl: 2, pr: 2, width: {lg: `calc(100% - ${drawerWidth}px)`}}}
      >
        {(fullyLoaded && !albums && (
          <Alert severity='info' sx={{mt: 3}}>Your account is empty, start to create new albums and upload your
            photos with the command line interface.</Alert>
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