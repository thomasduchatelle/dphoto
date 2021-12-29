import axios from "axios";
import React, {useEffect, useState} from 'react';
import './App.css';
import logo from './logo.svg';

interface Album {
  Name: String
  FolderName: String
  Start: Date
  End: Date
}

interface AlbumWithStats {
  Album: Album
  TotalCount: number
}

function App() {
  let [albums, setAlbums] = useState<AlbumWithStats[]>([])

  useEffect(() => {
    axios.get<AlbumWithStats[]>('/api/v1/albums').then(resp => setAlbums(resp.data))
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo"/>
        <p>
          Welcome to DPhoto !
        </p>
        {albums.length > 0 && (
          <div>
            <h4>{albums.length} album(s):</h4>
            <ul>
              {albums.map(a => (
                  <li>{a.Album.Name} starting on {a.Album.Start.toLocaleString()} [{a.TotalCount} media(s)]</li>
                )
              )}
            </ul>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
