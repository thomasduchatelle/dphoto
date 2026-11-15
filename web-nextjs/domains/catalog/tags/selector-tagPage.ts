import {TagMediaWithinADay, TagPageState} from "./language";

export interface TagPageSelection {
    tag: string
    medias: TagMediaWithinADay[]
    loaded: boolean
    isEmpty: boolean
    error?: Error
}

export const tagPageSelector = (state: TagPageState): TagPageSelection => ({
    tag: state.tag,
    medias: state.medias,
    loaded: state.loaded,
    isEmpty: state.loaded && !state.error && state.medias.length === 0,
    error: state.error,
});
