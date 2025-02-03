import {MediasLoadedAction, StartLoadingMediasAction, startLoadingMediasAction} from "./catalog-actions";
import {AlbumId} from "./catalog-state";
import {ActionObserver} from "./ActionObserver";
import {albumIdEquals} from "./utils-albumIdEquals";

export interface SelectAlbumQuery {
    loaded: boolean
    currentAlbumId: AlbumId | undefined
    albumId: AlbumId
}

export interface MediaPerDayLoaderInterface {
    loadMedias(albumId: AlbumId): Promise<MediasLoadedAction>
}

export class SelectAlbumHandler {

    constructor(
        private readonly dispatch: ActionObserver<MediasLoadedAction | StartLoadingMediasAction>,
        private readonly mediaPerDayLoader: MediaPerDayLoaderInterface,
        private readonly loadedAlbumId?: AlbumId,
        private readonly loadingAlbumId?: AlbumId,
    ) {
    }

    public onSelectAlbum(albumId: AlbumId) {
        if (this.loadedAlbumId && (
            (this.loadingAlbumId && !albumIdEquals(this.loadingAlbumId, albumId))
            || !albumIdEquals(this.loadedAlbumId, albumId)
        )) {

            this.dispatch(startLoadingMediasAction(albumId))
            return this.mediaPerDayLoader.loadMedias(albumId).then(this.dispatch)
        }

        return Promise.resolve()
    }
}

