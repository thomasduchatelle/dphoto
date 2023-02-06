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

export interface CatalogState {
    albumNotFound: boolean
    albums: Album[]
    selectedAlbum?: Album
    medias: MediaWithinADay[]
    error?: Error
    loadingMediasFor?: AlbumId
}

export interface CatalogStateV0 {
    fullyLoaded: boolean
    albumsLoaded: boolean
    albumNotFound: boolean
    albums: Album[]
    selectedAlbum?: Album
    medias: MediaWithinADay[]
}

export interface CatalogAPI {

    fetchAlbums(email: string): Promise<Album[]>

    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

// ---------------------------------
// TODO Delete beneath this line
// ---------------------------------


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

export function albumIdEquals(a: AlbumId, b: AlbumId): boolean {
    return a.owner === b.owner && a.folderName === b.folderName
}