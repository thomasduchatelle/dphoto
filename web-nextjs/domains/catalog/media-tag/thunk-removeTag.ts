import {ThunkDeclaration} from "@/libs/dthunks";
import {AlbumId, CatalogViewerState, MediaId} from "../language";
import {CatalogDispatch} from "../common/catalog-dispatch";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";
import {mediaTagAdded} from "./action-mediaTagAdded";
import {mediaTagRemoved} from "./action-mediaTagRemoved";

export interface RemoveTagPort {
    removeTag(owner: string, mediaId: MediaId, tag: string): Promise<void>;
}

export type RemoveTagThunk = (mediaId: MediaId, tag: string) => Promise<void>;

export async function removeTagThunk(
    dispatch: (action: ReturnType<typeof mediaTagAdded> | ReturnType<typeof mediaTagRemoved>) => void,
    port: RemoveTagPort,
    displayedAlbumId: AlbumId | undefined,
    mediaId: MediaId,
    tag: string
): Promise<void> {
    if (!displayedAlbumId) {
        return;
    }

    dispatch(mediaTagRemoved({mediaId, tag}));

    try {
        await port.removeTag(displayedAlbumId.owner, mediaId, tag);
    } catch {
        dispatch(mediaTagAdded({mediaId, tag}));
    }
}

export const removeTagDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { displayedAlbumId: AlbumId | undefined },
    RemoveTagThunk,
    CatalogDispatch & { adapter: RemoveTagPort }
> = {
    selector: (state: CatalogViewerState) => ({
        displayedAlbumId: getDisplayedAlbumId(state),
    }),
    factory: ({dispatch, adapter, partialState: {displayedAlbumId}}) => {
        return removeTagThunk.bind(null, dispatch, adapter, displayedAlbumId);
    },
};
