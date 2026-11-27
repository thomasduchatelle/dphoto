import {AlbumId, Media} from "../language";

export const FAVORITE_TAG = "favorite";

export interface Tag {
    tag: string
    namespace: string
    value: string
}

export interface TagMedia extends Media {
    albumId: AlbumId
    tags: string[]
}

export interface TagMediaWithinADay {
    day: Date
    medias: TagMedia[]
}

export interface TagPageState {
    tag: string
    medias: TagMediaWithinADay[]
    loaded: boolean
    error?: Error
}

export const initialTagPageState = (tag: string): TagPageState => ({
    tag,
    medias: [],
    loaded: false,
});
