import {AlbumId, MediaId} from "../language";
import {albumUrl} from "./album-url";

export function photoUrl(albumId: AlbumId, mediaId: MediaId): string {
    return `${albumUrl(albumId)}/photos/${encodeURIComponent(mediaId)}`;
}
