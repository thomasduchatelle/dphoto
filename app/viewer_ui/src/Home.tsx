import axios from "axios";
import React, {useEffect, useState} from 'react';
import './App.css';
// import logo from './logo.svg';
import {accessToken, owner} from "./security/google-authentication.service";
import {SecurityContextType} from "./security/security.model";
import {withSecurityContext} from "./security/with-security-context.hook";

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

const App = ({loggedUser}: SecurityContextType) => {
  const [albums, setAlbums] = useState<AlbumWithStats[]>([])

  useEffect(() => {
    axios.get<AlbumWithStats[]>(`/api/v1/owners/${owner}/albums`, {
      headers: {
        "Authorization": `Bearer ${accessToken}`,
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

export default withSecurityContext(App);
