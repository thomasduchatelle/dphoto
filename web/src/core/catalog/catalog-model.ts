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
    ownedBy?: OwnerDetails // only present when not owned by current user
    sharedWith: UserDetails[]
}

export interface OwnerDetails {
    name?: string
    users: UserDetails[]
}

export interface UserDetails {
    name: string
    email: string
    picture?: string
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

export function albumIdEquals(a?: AlbumId, b?: AlbumId): boolean {
    return !!a && a?.owner === b?.owner && a?.folderName === b?.folderName
}