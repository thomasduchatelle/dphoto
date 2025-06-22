import {Album, AlbumId, CatalogViewerState} from "../language";
import {AlbumsLoaded, albumsLoaded} from "../navigation";
import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {createAction} from "src/light-state-lib";

export interface CreateAlbumPayload {
    name: string
    start: Date
    end: Date
    forcedFolderName: string
}

export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumPayload): Promise<AlbumId>

    fetchAlbums(): Promise<Album[]>
}

export const createAlbum = createAction<
    { dispatch: (action: AlbumsLoaded) => void, createAlbumPort: CreateAlbumPort },
    CreateAlbumPayload,
    Promise<AlbumId>
>(
    "CreateAlbum",
    async ({dispatch, createAlbumPort}, request: CreateAlbumPayload) => {
        const albumId: AlbumId = await createAlbumPort.createAlbum(request);
        const albums: Album[] = await createAlbumPort.fetchAlbums();
        dispatch(albumsLoaded({albums, redirectTo: albumId}));
        return albumId;
    }
);

export type CreateAlbum = ReturnType<typeof createAlbum>;

export const createAlbumDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    CreateAlbum, // The thunk now directly returns the action
    CatalogFactoryArgs
> = {
    factory: ({dispatch, app}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        // The factory now returns the action creator itself, bound with its dependencies
        return createAlbum.bind(null, {dispatch, createAlbumPort: restAdapter});
    },
    selector: (_state: CatalogViewerState) => ({}),
};
