export interface AlbumId {
  owner: string
  folderName: string
}

export interface Album {
  albumId: AlbumId
  name: string
  start: Date
  end: Date
  totalCount: number
  temperature: number // number of media per day
  relativeTemperature: number
}

export enum MediaType {
  IMAGE,
  VIDEO,
  OTHER
}

export interface MediaWithinADay {
  day: Date
  medias: Media[]
}

export interface MediaId {
  hash: string
  size: number
}

export interface Media {
  id: string
  type: MediaType
  time: Date
  uiRelativePath: string // uiRelativePath is the WEB UI internal link (from the album)
  contentPath: string
  source: string
}

export interface WebAdapter {
  redirectToAlbum(albumId: AlbumId): void;

  renderNoAlbums(): void;

  renderAlbumNotPresent(albums: Album[], albumId: AlbumId): void;

  renderAlbumsWithMedia(albums: Album[], selectedAlbum: Album, medias: MediaWithinADay[]): void;

  renderAlbumsList(albums: Album[]): void;
}

export interface AlbumsLogicCache {
  owner: string
  albums: Album[]
}

