import {toggleFavoriteDeclaration} from "./thunk-toggleFavorite";

export type {ToggleFavoritePort, ToggleFavoriteThunk} from "./thunk-toggleFavorite";
export {mediaTagsUpdated} from "./action-mediaTagsUpdated";

/**
 * Thunks related to tagging a media.
 *
 * Expected handler types:
 * - `toggleFavorite`: `(mediaId: MediaId, currentTags: string[]) => Promise<void>`
 */
export const mediaTagThunks = {
    toggleFavorite: toggleFavoriteDeclaration,
};
