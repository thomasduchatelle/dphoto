import axios from "axios";
import React, {useEffect, useState} from 'react';
import {getAppContext, useAuthenticatedUser} from "./core/application";

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
    axios.get<AlbumWithStats[]>(`/api/v1/owners/${loggedUser?.email}/albums`, {
      headers: {
        "Authorization": `Bearer ${getAppContext().oauthService.getAccessToken()}`,
      }
    }).then(resp => setAlbums(resp.data))
  }, [])

  return (
    <header className="App-header">
      {/*<img src={logo} className="App-logo" alt="logo"/>*/}
      <p>
        Welcome to DPhoto{loggedUser && `, ${loggedUser.name}`} !
      </p>
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
    </header>
  );
}

export default App;
