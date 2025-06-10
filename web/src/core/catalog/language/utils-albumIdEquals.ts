import {AlbumId} from "./catalog-state";

export function albumIdEquals(a?: AlbumId, b?: AlbumId): boolean {
    return !!a && a?.owner === b?.owner && a?.folderName === b?.folderName
}