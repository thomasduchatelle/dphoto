import {createAction} from "@/libs/daction";
import {TagPageState} from "./language";

interface TagMediasLoadFailedPayload {
    tag: string
    error: Error
}

export const tagMediasLoadFailed = createAction<TagPageState, TagMediasLoadFailedPayload>(
    "tagMediasLoadFailed",
    (current: TagPageState, {tag, error}: TagMediasLoadFailedPayload): TagPageState => {
        if (current.tag !== tag) {
            return current;
        }

        return {
            ...current,
            medias: [],
            loaded: true,
            error,
        };
    }
);

export type TagMediasLoadFailed = ReturnType<typeof tagMediasLoadFailed>;
