import {albumCreateThunks} from "./album-create";
import {navigationThunks} from "./navigation";
import {albumDeleteThunks} from "./album-delete";
import {albumEditDatesThunks} from "./album-edit-dates";
import {albumEditNameThunks} from "./album-edit-name";
import {ThunkDeclaration} from "@/libs/dthunks";
import {sharingThunks} from "./sharing";
import {baseNameEditThunks} from "./base-name-edit";

export * from "./common/catalog-factory-args";
export type {RevokeAlbumAccessAPI} from "./sharing/thunk-revokeAlbumAccess";
export type {GrantAlbumAccessAPI} from "./sharing/thunk-grantAlbumAccess";
export type {DeleteAlbumThunk, DeleteAlbumPort} from "./album-delete/thunk-deleteAlbum";
export type {SaveAlbumNamePort} from "./album-edit-name/thunk-saveAlbumName";


export const catalogThunks = {
    ...navigationThunks,
    ...sharingThunks,
    ...albumCreateThunks,
    ...albumDeleteThunks,
    ...albumEditDatesThunks,
    ...albumEditNameThunks,
    ...baseNameEditThunks,
};

// Dynamically infer the interface from catalogThunks
export type CatalogThunksInterface = {
    [K in keyof typeof catalogThunks]:
    typeof catalogThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};