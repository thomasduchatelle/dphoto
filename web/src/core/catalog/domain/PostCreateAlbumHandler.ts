import {AlbumId} from "./catalog-state";
import {CreateAlbumListener} from "./CreateAlbumController";
import {AlbumsLoadedAction, catalogActions} from "./catalog-reducer-v2";
import {FetchAlbumsPort} from "../thunks";

export class PostCreateAlbumHandler implements CreateAlbumListener {

    constructor(
        private readonly dispatch: (action: AlbumsLoadedAction) => void,
        private readonly fetchAlbumsPort: FetchAlbumsPort,
    ) {
    }

    onAlbumCreated = async (albumId: AlbumId): Promise<void> => {
        const albums = await this.fetchAlbumsPort.fetchAlbums()
        this.dispatch(catalogActions.albumsLoadedAction({albums: albums, redirectTo: albumId}))
    }
}
