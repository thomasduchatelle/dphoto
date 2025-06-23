import {closeDeleteAlbumDialogDeclaration} from "./thunk-closeDeleteAlbumDialog";
import {openDeleteAlbumDialogDeclaration} from "./thunk-openDeleteAlbumDialog";
import {deleteAlbumDeclaration} from "./thunk-deleteAlbum";

export * from "./selector-deleteDialogSelector"

/**
 * Thunks related to album deletion.
 *
 * Expected handler types:
 * - `closeDeleteAlbumDialog`: `() => void`
 * - `openDeleteAlbumDialog`: `() => void`
 * - `deleteAlbum`: `(albumIdToDelete: AlbumId) => Promise<void>`
 */
export const albumDeleteThunks = {
    closeDeleteAlbumDialog: closeDeleteAlbumDialogDeclaration,
    openDeleteAlbumDialog: openDeleteAlbumDialogDeclaration,
    deleteAlbum: deleteAlbumDeclaration,
};
