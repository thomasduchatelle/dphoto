import {Container, Paper, Typography} from "@mui/material";
import React, {useEffect, useState} from 'react';
import {authenticatedAxios, useAuthenticatedUser} from "../../core/application";

interface Album {
  Name: string
  FolderName: string
  Start: Date
  End: Date
}

interface AlbumWithStats {
  Album: Album
  TotalCount: number
}

const App = () => {
  const loggedUser = useAuthenticatedUser()
  const [albums, setAlbums] = useState<AlbumWithStats[]>([])

  useEffect(() => {
    authenticatedAxios().get<AlbumWithStats[]>(`/api/v1/owners/${loggedUser?.email}/albums`).then(resp => setAlbums(resp.data))
  }, [])

  return (
    <Container maxWidth={false}>
      <Paper sx={{p: 2, mt: 2}}>
        <Typography component='h1' variant='h3'>
          Welcome to DPhoto{loggedUser && `, ${loggedUser.name}`} !
        </Typography>
        {albums.length > 0 && (
          <div>
            <h4>{albums.length} album(s):</h4>
            <ul>
              {albums.map(a => (
                  <li key={a.Album.FolderName}>{a.Album.Name} starting
                    on {a.Album.Start.toLocaleString()} [{a.TotalCount} media(s)]</li>
                )
              )}
            </ul>
          </div>
        )}
      </Paper>
    </Container>
  );
}

export default App;
