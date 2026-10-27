import {AlbumId, MediaId} from "../language";
import {albumUrl} from "./album-url";

export function photoUrl(albumId: AlbumId, mediaId: MediaId, options?: {autoplay: boolean}): string {
    const url = `${albumUrl(albumId)}/photos/${encodeURIComponent(mediaId)}`;
    return options?.autoplay ? `${url}?play` : url;
}
