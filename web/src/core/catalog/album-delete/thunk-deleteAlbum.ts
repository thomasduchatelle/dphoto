import type {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {isDeleteAlbumError} from "../adapters/api";
import {Album, AlbumId, albumIdEquals, CatalogViewerState} from "../language";
import {deleteAlbumStarted} from "./action-deleteAlbumStarted";
import {albumDeleteFailed} from "./action-albumDeleteFailed";
import {albumDeleted} from "./action-albumDeleted";
import {ThunkDeclaration} from "src/libs/thunks";


export interface DeleteAlbumPort {
    deleteAlbum(albumId: AlbumId): Promise<void>;

    fetchAlbums(): Promise<Album[]>;
}

export type DeleteAlbumThunk = (albumIdToDelete: AlbumId) => Promise<void>;

function getSelectedAlbumId(state: CatalogViewerState): AlbumId | undefined {
    return state.loadingMediasFor ?? state.mediasLoadedFromAlbumId;
}

function errorToMessage(error: any): string {
    if (isDeleteAlbumError(error)) {
        console.error(`Client error [${error.code}]: ${error.message}`);
        return error.message
    }

    console.log(`Unexpected error: ${error}`);
    return "A technical error prevented the album to be deleted, please report it to a developer.";
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
        dispatch(albumDeleteFailed(errorToMessage(error)));
        return;
    }

    try {
        const albums = await port.fetchAlbums();
        const redirectTo = albumIdEquals(selectedAlbumId, albumIdToDelete) && !!albums ? albums[0]?.albumId : undefined; // Added optional chaining
        dispatch(albumDeleted({albums, redirectTo}));

    } catch (error) {
        dispatch(albumDeleteFailed(`Failed to fetch albums after deletion: ${error}`));
    }
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
