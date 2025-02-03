import {AlbumId} from "./catalog-state";
import {CreateAlbumListener} from "./CreateAlbumController";
import {AlbumsLoadedAction} from "./catalog-actions";
import {FetchAlbumsPort} from "./CatalogLoader";

export class PostCreateAlbumHandler implements CreateAlbumListener {

    constructor(
        private readonly dispatch: (action: AlbumsLoadedAction) => void,
        private readonly fetchAlbumsPort: FetchAlbumsPort,
    ) {
    }

    onAlbumCreated = async (albumId: AlbumId): Promise<void> => {
        const albums = await this.fetchAlbumsPort.fetchAlbums()
        this.dispatch({
            type: 'AlbumsLoadedAction',
            albums: albums,
            redirectTo: albumId,
        })
    }
}