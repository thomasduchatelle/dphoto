import {closeSharingModalDeclaration, grantAlbumAccessDeclaration, openSharingModalDeclaration, revokeAlbumAccessDeclaration} from "./sharing";
import {createAlbumDeclaration} from "./album-create";
import {onAlbumFilterChangeDeclaration, onPageRefreshDeclaration} from "./navigation";
import {closeDeleteAlbumDialogDeclaration, deleteAlbumDeclaration, openDeleteAlbumDialogDeclaration} from "./album-delete";
import {
    albumEditDatesThunks,
    closeEditDatesDialogDeclaration,
    openEditDatesDialogDeclaration,
    updateAlbumDatesDeclaration,
    updateEditDatesDialogEndDateDeclaration,
    updateEditDatesDialogStartDateDeclaration
} from "./album-edit-dates";
import {ThunkDeclaration} from "src/libs/dthunks";

export * from "./common/catalog-factory-args";
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
    ...albumEditDatesThunks, // Aggregate albumEditDatesThunks here
};

// Dynamically infer the interface from catalogThunks
export type CatalogThunksInterface = {
    [K in keyof typeof catalogThunks]:
    typeof catalogThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};
