import type {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {Album, AlbumId, CatalogViewerState, getErrorMessage, Media} from "../language";
import {deleteAlbumStarted} from "./action-deleteAlbumStarted";
import {albumDeleteFailed} from "./action-albumDeleteFailed";
import {ThunkDeclaration} from "@/libs/dthunks";
import {loadAlbumsAndMedias} from "../navigation/utils-loadAlbumsAndMedias";


export interface DeleteAlbumPort {
    deleteAlbum(albumId: AlbumId): Promise<void>;

    fetchAlbums(): Promise<Album[]>;

    fetchMedias(albumId: AlbumId): Promise<Media[]>;
}

export type DeleteAlbumThunk = (albumIdToDelete: AlbumId) => Promise<void>;

function getSelectedAlbumId(state: CatalogViewerState): AlbumId | undefined {
    return state.loadingMediasFor ?? state.mediasLoadedFromAlbumId;
}

export async function deleteAlbumThunk(
    dispatch: (action: any) => void,
    port: DeleteAlbumPort,
    selectedAlbumId: AlbumId | undefined,
    albumIdToDelete: AlbumId
): Promise<void> {
    dispatch(deleteAlbumStarted());

    try {
        await port.deleteAlbum(albumIdToDelete);
    } catch (error) {
        dispatch(albumDeleteFailed(getErrorMessage(error) ?? "A technical error prevented the album to be deleted, please report it to a developer."));
        return;
    }

    const action = await loadAlbumsAndMedias(port, selectedAlbumId)
    dispatch(action);
}

export const deleteAlbumDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { selectedAlbumId: AlbumId | undefined },
    DeleteAlbumThunk,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({
        selectedAlbumId: getSelectedAlbumId(state)
    }),
    factory: ({dispatch, app, partialState: {selectedAlbumId}}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        return deleteAlbumThunk.bind(null, dispatch, restAdapter, selectedAlbumId);
    }
};
