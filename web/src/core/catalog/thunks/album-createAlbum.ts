import {Album, AlbumId, AlbumsLoadedAction, catalogActions, CatalogViewerState} from "../domain";
import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogFactoryArgs} from "./catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";

export interface CreateAlbumRequest {
    name: string
    start: Date
    end: Date
    forcedFolderName: string
}

export type CreateAlbumThunk = (request: CreateAlbumRequest) => Promise<AlbumId>


export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumRequest): Promise<AlbumId>

    fetchAlbums(): Promise<Album[]>
}

export async function createAlbumThunk(
    dispatch: (action: AlbumsLoadedAction) => void,
    createAlbumPort: CreateAlbumPort,
    request: CreateAlbumRequest
): Promise<AlbumId> {
    const albumId: AlbumId = await createAlbumPort.createAlbum(request);
    const albums: Album[] = await createAlbumPort.fetchAlbums();
    dispatch(catalogActions.albumsLoadedAction({albums, redirectTo: albumId}));

    return albumId
}

export const createAlbumDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    CreateAlbumThunk,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        return createAlbumThunk.bind(null, dispatch, restAdapter);
    },
    selector: (_state: CatalogViewerState) => ({}),
};
