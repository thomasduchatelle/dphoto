export * from "./action-albumDeleted"
export * from "./action-albumDeleteFailed"
export * from "./action-deleteAlbumDialogClosed"
export * from "./action-deleteAlbumDialogOpened"
export * from "./action-deleteAlbumStarted"
export * from "./selector-deleteDialogSelector"
export * from "./thunk-deleteAlbum"

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
