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
    sharedWith: Sharing[]
}

export interface OwnerDetails {
    name?: string
    users: UserDetails[]
}

export enum SharingType {
    visitor = "visitor",
    contributor = "contributor",
}

export interface Sharing {
    user: UserDetails
    role: SharingType
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

export type MediaId = string

export interface Media {
    id: MediaId
    type: MediaType
    time: Date
    uiRelativePath: string // uiRelativePath is the WEB UI internal link (from the album)
    contentPath: string
    source: string
}

export function albumIdEquals(a?: AlbumId, b?: AlbumId): boolean {
    return !!a && a?.owner === b?.owner && a?.folderName === b?.folderName
}