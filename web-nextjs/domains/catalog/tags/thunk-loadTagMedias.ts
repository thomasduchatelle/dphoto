import {Action} from "@/libs/daction";
import {ThunkDeclaration} from "@/libs/dthunks";
import {TagMedia, TagPageState} from "./language";
import {tagMediasLoaded, TagMediasLoaded} from "./action-tagMediasLoaded";
import {tagMediasLoadFailed, TagMediasLoadFailed} from "./action-tagMediasLoadFailed";

export interface LoadTagMediasPort {
    fetchMediasByTag(tag: string): Promise<TagMedia[]>
}

export async function loadTagMedias(
    dispatch: (action: TagMediasLoaded | TagMediasLoadFailed) => void,
    port: LoadTagMediasPort,
    tag: string,
): Promise<void> {
    try {
        const medias = await port.fetchMediasByTag(tag);
        dispatch(tagMediasLoaded({tag, medias}));
    } catch (error) {
        dispatch(tagMediasLoadFailed({
            tag,
            error: new Error(`failed to load medias of tag '${tag}'`, {cause: error}),
        }));
    }
}

export interface TagPageDispatch {
    dispatch: (action: Action<TagPageState, any>) => void;
}

export const loadTagMediasDeclaration: ThunkDeclaration<
    TagPageState,
    Record<string, never>,
    (tag: string) => Promise<void>,
    TagPageDispatch & { adapter: LoadTagMediasPort }
> = {
    selector: () => ({}),
    factory: ({dispatch, adapter}) => loadTagMedias.bind(null, dispatch, adapter),
};
