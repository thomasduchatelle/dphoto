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

export interface Media {
  id: string
  type: MediaType
  time: Date
  path: string
  source: string
}

export interface WebAdapter {
  redirectToAlbum(albumId: AlbumId): void;

  renderNoAlbums(): void;

  renderAlbumNotPresent(albums: Album[]): void;

  renderAlbumsWithMedia(albums: Album[], selectedAlbum: Album, medias: Media[]): void;
}

export interface CacheAdapter {
  getCachedAlbums(): [string, Album[]]

  cacheAlbums(owner: string, albums: Album[]): void
}
