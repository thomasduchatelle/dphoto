import {Alert, Box, Drawer, Toolbar} from "@mui/material";
import React, {useEffect, useMemo, useRef, useState} from 'react';
import {useNavigate, useParams} from "react-router-dom"
import AppNavComponent from "../../../../components/app-nav.component";
import UserMenu from "../../../../components/user.menu";
import {useMustBeAuthenticated} from "../../../../core/application";
import {Album, AlbumId, AlbumsPageCase, Media} from "../../../../core/domain/catalog";
import AlbumsListComponent from "../albums-list.component";
import {MediaListComponent} from "../media-list.component";

type AlbumPageUrlParams = {
  owner: string | undefined,
  album: string | undefined,
}

interface AlbumPageState {
  fullyLoaded: boolean
  albumNotFound: boolean
  albums: Album[]
  selectedAlbum?: Album
  medias: Media[]
}

interface AlbumPageCache {
  owner: string
  albums: Album[]
}

const AlbumPage = () => {
  const [loggedUser, axiosInstance, signOutCase] = useMustBeAuthenticated()
  const {owner, album} = useParams<AlbumPageUrlParams>()
  const [state, setState] = useState<AlbumPageState>({fullyLoaded: false, albumNotFound: false, albums: [], medias: []})
  const cache = useRef<AlbumPageCache>({owner: '', albums: []})
  const navigate = useNavigate()

  const engine = useMemo(() => new AlbumsPageCase(axiosInstance,
    {
      cacheAlbums(owner: string, albums: Album[]): void {
        cache.current = {owner, albums}
      },
      getCachedAlbums(): [string, Album[]] {
        return [cache.current.owner, cache.current.albums];
      }
    },
    {
      redirectToAlbum(albumId: AlbumId): void {
        navigate(`/albums/${albumId.owner}/${albumId.folderName}`, {replace: false})
      },
      renderAlbumNotPresent(albums: Album[]): void {
        setState(current => ({...current, albums, fullyLoaded: true, albumNotFound: true, medias: []}))
      },
      renderAlbumsWithMedia(albums: Album[], selectedAlbum: Album, medias: Media[]): void {
        setState(current => ({...current, albums, medias, selectedAlbum, fullyLoaded: true, albumNotFound: false}))
      },
      renderNoAlbums(): void {
        setState(current => ({...current, fullyLoaded: true, albums: [], medias: [], albumNotFound: false}))
      }
    }), [axiosInstance, setState, navigate])

  useEffect(() => {
    if (!owner || !album) {
      engine.redirectToDefaultAlbum(loggedUser).catch(err => console.log(`Error: ${err}\n${err.stack}`))
    } else {
      engine.refreshPage(loggedUser, {
        owner,
        folderName: album
      }).catch(err => console.log(`Error: ${err}\n${err.stack}`))
    }
  }, [engine, loggedUser, owner, album])

  const drawerWidth = 450
  return (
    <Box sx={{display: 'flex'}}>
      <AppNavComponent
        rightContent={<UserMenu user={loggedUser} onLogout={signOutCase.logout}/>}
      />
      <Box
        component="nav"
        sx={{width: {sm: drawerWidth}, flexShrink: {sm: 0}}}
        aria-label="mailbox folders"
      >
        <Drawer
          variant="permanent"
          sx={{
            width: drawerWidth,
            display: {xs: 'none', sm: 'block'},
            flexShrink: 0,
            [`& .MuiDrawer-paper`]: {width: drawerWidth, boxSizing: 'border-box'},
          }}
        >
          <Toolbar/>
          <AlbumsListComponent albums={state.albums} loaded={state.fullyLoaded} selected={state.selectedAlbum}
                               onSelect={engine.selectAlbum}/>
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{flexGrow: 1, pl: 2, pr: 2, width: {sm: `calc(100% - ${drawerWidth}px)`}}}
      >
        <Toolbar/>
        {(state.fullyLoaded && state.albumNotFound && (
          <Alert severity='info' sx={{mt: 3}}>Your account is empty, start to create new albums and upload your
            photos with the command line interface.</Alert>
        )) || (
          <MediaListComponent medias={state.medias} loaded={state.fullyLoaded} albumNotFound={state.albumNotFound}/>
        )}
      </Box>
    </Box>
  );
}

export default AlbumPage;
