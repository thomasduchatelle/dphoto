import {AlbumId} from "../language";

export function albumUrl(albumId: AlbumId): string {
    return `/albums/${encodeURIComponent(albumId.owner)}/${encodeURIComponent(albumId.folderName)}`;
}
