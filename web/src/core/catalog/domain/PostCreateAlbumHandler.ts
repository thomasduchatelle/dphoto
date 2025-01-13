import {AlbumId} from "./catalog-state";
import {CreateAlbumListener} from "./CreateAlbumController";
import {AlbumsAndMediasLoadedAction, MediaFailedToLoadAction, NoAlbumAvailableAction} from "./catalog-actions";
import {CatalogViewerLoader} from "./CatalogViewerLoader";


export class PostCreateAlbumHandler implements CreateAlbumListener {

    constructor(
        private readonly dispatch: (action: MediaFailedToLoadAction | AlbumsAndMediasLoadedAction | NoAlbumAvailableAction) => void,
        private readonly catalogViewerLoader: CatalogViewerLoader,
    ) {
    }


    onAlbumCreated = (albumId: AlbumId): Promise<void> => {
        return this.catalogViewerLoader.loadInitialCatalog({
            albumId,
        }).then(this.dispatch)
    }

}