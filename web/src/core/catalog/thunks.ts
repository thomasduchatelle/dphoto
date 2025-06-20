import {closeSharingModalDeclaration, grantAlbumAccessDeclaration, openSharingModalDeclaration, revokeAlbumAccessDeclaration} from "./sharing";
import type {ThunkDeclaration} from "../thunk-engine";
import {createAlbumDeclaration} from "./album-create";
import {onAlbumFilterChangeDeclaration, onPageRefreshDeclaration} from "./navigation";
import {closeDeleteAlbumDialogDeclaration, deleteAlbumDeclaration, openDeleteAlbumDialogDeclaration} from "./album-delete";
import {closeEditDatesDialogDeclaration, openEditDatesDialogDeclaration} from "./album-edit-dates";

export * from "./common/catalog-factory-args";
export type {FetchAlbumsPort} from "./navigation/thunk-onPageRefresh";
export type {RevokeAlbumAccessAPI} from "./sharing/thunk-revokeAlbumAccess";
export type {GrantAlbumAccessAPI} from "./sharing/thunk-grantAlbumAccess";
export type {CreateAlbumThunk, CreateAlbumRequest, CreateAlbumPort} from "./album-create/album-createAlbum";
export type {DeleteAlbumThunk, DeleteAlbumPort} from "./album-delete/thunk-deleteAlbum";
export {openDeleteAlbumDialogThunk} from "./album-delete/thunk-openDeleteAlbumDialog";
export {closeDeleteAlbumDialogThunk} from "./album-delete/thunk-closeDeleteAlbumDialog";

export const catalogThunks = {
    onPageRefresh: onPageRefreshDeclaration,
    onAlbumFilterChange: onAlbumFilterChangeDeclaration,
    openSharingModal: openSharingModalDeclaration,
    closeSharingModal: closeSharingModalDeclaration,
    revokeAlbumAccess: revokeAlbumAccessDeclaration,
    grantAlbumSharing: grantAlbumAccessDeclaration,
    createAlbum: createAlbumDeclaration,
    deleteAlbum: deleteAlbumDeclaration,
    openDeleteAlbumDialog: openDeleteAlbumDialogDeclaration,
    closeDeleteAlbumDialog: closeDeleteAlbumDialogDeclaration,
    openEditDatesDialog: openEditDatesDialogDeclaration,
    closeEditDatesDialog: closeEditDatesDialogDeclaration,
};

// Dynamically infer the interface from catalogThunks
export type CatalogThunksInterface = {
    [K in keyof typeof catalogThunks]:
    typeof catalogThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};
