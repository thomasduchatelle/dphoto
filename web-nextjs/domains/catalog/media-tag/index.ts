import {addTagDeclaration} from "./thunk-addTag";
import {removeTagDeclaration} from "./thunk-removeTag";

export type {AddTagPort, AddTagThunk} from "./thunk-addTag";
export type {RemoveTagPort, RemoveTagThunk} from "./thunk-removeTag";
export {mediaTagAdded} from "./action-mediaTagAdded";
export {mediaTagRemoved} from "./action-mediaTagRemoved";

/**
 * Thunks related to tagging a media.
 *
 * Expected handler types:
 * - `addTag`: `(mediaId: MediaId, tag: string) => Promise<void>`
 * - `removeTag`: `(mediaId: MediaId, tag: string) => Promise<void>`
 */
export const mediaTagThunks = {
    addTag: addTagDeclaration,
    removeTag: removeTagDeclaration,
};
