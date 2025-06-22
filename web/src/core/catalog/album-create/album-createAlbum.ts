import {Album, AlbumId, CatalogViewerState} from "../language";
import {AlbumsLoaded, albumsLoaded} from "../navigation";
import type {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {ThunkDeclaration} from "src/libs/thunks";

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
    dispatch: (action: AlbumsLoaded) => void,
    createAlbumPort: CreateAlbumPort,
    request: CreateAlbumRequest
): Promise<AlbumId> {
    const albumId: AlbumId = await createAlbumPort.createAlbum(request);
    const albums: Album[] = await createAlbumPort.fetchAlbums();
    dispatch(albumsLoaded({albums, redirectTo: albumId}));

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
