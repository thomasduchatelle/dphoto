import {LoadTagMediasPort} from "../thunk-loadTagMedias";
import {MockTagsAdapter} from "./mock-tags-adapter";

// MOCK adapter against the pinned list-medias-by-tag contract (design.md).
// Swapped for a REST adapter once the 02.1 backend endpoint is available.
export function newTagsAdapter(): LoadTagMediasPort {
    return new MockTagsAdapter();
}
