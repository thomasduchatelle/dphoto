import {onPageRefreshDeclaration} from "./thunks-onPageRefresh";
import {onAlbumFilterChangeDeclaration} from "./thunks-onAlbumFilterChange";
import {openSharingModalDeclaration} from "./share-openSharingModal";
import {closeSharingModalDeclaration} from "./share-closeSharingModal";
import {revokeAlbumSharingDeclaration} from "./share-revokeAlbumSharing";
import {grantAlbumSharingDeclaration} from "./share-grantAlbumSharing";
import type {ThunkDeclaration} from "../../thunk-engine";
import {createAlbumDeclaration} from "./album-createAlbum";
import {openDeleteAlbumDialogDeclaration} from "./album-openDeleteAlbumDialog";
import {closeDeleteAlbumDialogDeclaration} from "./album-closeDeleteAlbumDialog";
import {deleteAlbumDeclaration} from "./album-deleteAlbum";

export * from "./catalog-factory-args";
export type {FetchAlbumsPort} from "./thunks-onPageRefresh";
export type {revokeAlbumSharingAPI} from "./share-revokeAlbumSharing";
export type {GrantAlbumSharingAPI} from "./share-grantAlbumSharing";
export type {CreateAlbumThunk, CreateAlbumRequest, CreateAlbumPort} from "./album-createAlbum";
export type {DeleteAlbumThunk, DeleteAlbumPort} from "./album-deleteAlbum";
export {openDeleteAlbumDialogThunk} from "./album-openDeleteAlbumDialog";
export {closeDeleteAlbumDialogThunk} from "./album-closeDeleteAlbumDialog";

export const catalogThunks = {
    onPageRefresh: onPageRefreshDeclaration,
    onAlbumFilterChange: onAlbumFilterChangeDeclaration,
    openSharingModal: openSharingModalDeclaration,
    closeSharingModal: closeSharingModalDeclaration,
    revokeAlbumSharing: revokeAlbumSharingDeclaration,
    grantAlbumSharing: grantAlbumSharingDeclaration,
    createAlbum: createAlbumDeclaration,
    deleteAlbum: deleteAlbumDeclaration,
    openDeleteAlbumDialog: openDeleteAlbumDialogDeclaration,
    closeDeleteAlbumDialog: closeDeleteAlbumDialogDeclaration,
};

// Dynamically infer the interface from catalogThunks
export type CatalogThunksInterface = {
    [K in keyof typeof catalogThunks]:
    typeof catalogThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};
