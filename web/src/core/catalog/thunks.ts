import {albumCreateThunks} from "./album-create";
import {navigationThunks} from "./navigation";
import {albumDeleteThunks} from "./album-delete";
import {albumEditDatesThunks} from "./album-edit-dates";
import {ThunkDeclaration} from "src/libs/dthunks";
import {sharingThunks} from "./sharing";

export * from "./common/catalog-factory-args";
export type {RevokeAlbumAccessAPI} from "./sharing/thunk-revokeAlbumAccess";
export type {GrantAlbumAccessAPI} from "./sharing/thunk-grantAlbumAccess";
export type {CreateAlbumThunk, CreateAlbumRequest, CreateAlbumPort} from "./album-create/album-createAlbum";
export type {DeleteAlbumThunk, DeleteAlbumPort} from "./album-delete/thunk-deleteAlbum";


export const catalogThunks = {
    ...navigationThunks,
    ...sharingThunks,
    ...albumCreateThunks,
    ...albumDeleteThunks,
    ...albumEditDatesThunks,
};

// Dynamically infer the interface from catalogThunks
export type CatalogThunksInterface = {
    [K in keyof typeof catalogThunks]:
    typeof catalogThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};
