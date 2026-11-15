import {ThunkDeclaration} from "@/libs/dthunks";
import {loadTagMediasDeclaration} from "./thunk-loadTagMedias";

export type {LoadTagMediasPort} from "./thunk-loadTagMedias";

export const tagPageThunks = {
    loadTagMedias: loadTagMediasDeclaration,
};

export type TagPageThunksInterface = {
    [K in keyof typeof tagPageThunks]:
    typeof tagPageThunks[K] extends ThunkDeclaration<any, any, infer F, any> ? F : never
};
