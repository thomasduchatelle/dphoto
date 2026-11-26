import {ThunkDeclaration} from "@/libs/dthunks";
import {AlbumId, CatalogViewerState, MediaId} from "../language";
import {CatalogDispatch} from "../common/catalog-dispatch";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";
import {mediaTagAdded} from "./action-mediaTagAdded";
import {mediaTagRemoved} from "./action-mediaTagRemoved";

export interface AddTagPort {
    addTag(owner: string, mediaId: MediaId, tag: string): Promise<void>;
}

export type AddTagThunk = (mediaId: MediaId, tag: string) => Promise<void>;

export async function addTagThunk(
    dispatch: (action: ReturnType<typeof mediaTagAdded> | ReturnType<typeof mediaTagRemoved>) => void,
    port: AddTagPort,
    displayedAlbumId: AlbumId | undefined,
    mediaId: MediaId,
    tag: string
): Promise<void> {
    if (!displayedAlbumId) {
        return;
    }

    dispatch(mediaTagAdded({mediaId, tag}));

    try {
        await port.addTag(displayedAlbumId.owner, mediaId, tag);
    } catch {
        dispatch(mediaTagRemoved({mediaId, tag}));
    }
}

export const addTagDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { displayedAlbumId: AlbumId | undefined },
    AddTagThunk,
    CatalogDispatch & { adapter: AddTagPort }
> = {
    selector: (state: CatalogViewerState) => ({
        displayedAlbumId: getDisplayedAlbumId(state),
    }),
    factory: ({dispatch, adapter, partialState: {displayedAlbumId}}) => {
        return addTagThunk.bind(null, dispatch, adapter, displayedAlbumId);
    },
};
