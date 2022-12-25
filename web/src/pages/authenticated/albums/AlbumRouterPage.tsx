import {Box, Toolbar} from "@mui/material";
import React, {useEffect, useMemo, useRef, useState} from 'react';
import {useNavigate, useParams} from "react-router-dom"
import AppNavComponent from "../../../components/AppNav";
import UserMenu from "../../../components/user.menu";
import {useMustBeAuthenticated} from "../../../core/application";
import useWindowDimensions from "../../../core/window-utils";
import AlbumsList from "./AlbumsList";
import {Album, AlbumId, AlbumsLogic, AlbumsLogicCache, MediaWithinADay} from "./logic";
import MediasPage from "./MediasPage";
import MobileNavigation from "./MobileNavigation";

type AlbumPageUrlParams = {
  owner: string | undefined,
  album: string | undefined,
}

interface AlbumPageState {
  fullyLoaded: boolean
  albumsLoaded: boolean
  albumNotFound: boolean
  albums: Album[]
  selectedAlbum?: Album
  medias: MediaWithinADay[]
}

export default function AlbumRouterPage() {
  const mustBeAuthenticated = useMustBeAuthenticated()
  const {owner, album} = useParams<AlbumPageUrlParams>()
  const [state, setState] = useState<AlbumPageState>({
    fullyLoaded: false,
    albumsLoaded: false,
    albumNotFound: false,
    albums: [],
    medias: [],
  })
  const cache = useRef<AlbumsLogicCache>({owner: '', albums: []})
  const navigate = useNavigate()
  const {width} = useWindowDimensions()

  const engine = useMemo(() => new AlbumsLogic(mustBeAuthenticated, cache, width,
    {
      redirectToAlbum(albumId: AlbumId): void {
        navigate(`/albums/${albumId.owner}/${albumId.folderName}`, {replace: false})
      },
      renderAlbumNotPresent(albums: Album[], albumId: AlbumId): void {
        setState({
          albums,
          fullyLoaded: true,
          albumsLoaded: true,
          albumNotFound: true,
          medias: [],
          selectedAlbum: {
            albumId: albumId,
            end: new Date(),
            name: "not found",
            relativeTemperature: 0,
            start: new Date(),
            temperature: 0,
            totalCount: 0
          },
        })
      },
      renderAlbumsList(albums: Album[]): void {
        setState(current => ({...current, albums, albumsLoaded: true, fullyLoaded: false}))
      },
      renderAlbumsWithMedia(albums: Album[], selectedAlbum: Album, medias: MediaWithinADay[]): void {
        setState(current => ({
          ...current,
          albums,
          medias,
          selectedAlbum,
          fullyLoaded: true,
          albumsLoaded: true,
          albumNotFound: false,
        }))
      },
      renderNoAlbums(): void {
        setState(current => ({
          ...current,
          fullyLoaded: true,
          albumsLoaded: true,
          albums: [],
          medias: [],
          albumNotFound: false,
        }))
      }
    }), [mustBeAuthenticated, cache, width, setState, navigate]) //mustBeAuthenticated

  useEffect(() => {
    if (!owner || !album) {
      engine.loadAlbumsPage().catch(err => console.log(`Error: ${err}\n${err.stack}`))
    } else {
      engine.loadMediasPage({
        owner,
        folderName: album
      }).catch(err => console.log(`Error: ${err}\n${err.stack}`))
    }
  }, [engine, owner, album])

  return (
    <Box>
      <AppNavComponent
        rightContent={<UserMenu user={mustBeAuthenticated.loggedUser} onLogout={mustBeAuthenticated.signOutCase.logout}/>}
      />
      <Toolbar/>
      <Box sx={{mt: 2, pl: 2, pr: 2, display: {lg: 'none'}}}>
        <MobileNavigation album={owner && album ? state.selectedAlbum : undefined}/>
      </Box>
      {!owner || !album ? (
        <AlbumsList albums={state.albums}
                    loaded={state.albumsLoaded}
                    selected={state.selectedAlbum}
                    onSelect={engine.selectAlbum}/>
      ) : (
        <MediasPage
          albums={state.albums}
          albumNotFound={state.albumNotFound}
          fullyLoaded={state.fullyLoaded}
          medias={state.medias}
          selectAlbum={engine.selectAlbum}
          selectedAlbum={state.selectedAlbum}
        />
      )}
    </Box>
  );
}
