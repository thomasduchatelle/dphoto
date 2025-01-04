import {MediasLoadedAction, StartLoadingMediasAction, startLoadingMediasAction} from "./catalog-actions";
import {AlbumId} from "./catalog-state";
import {albumIdEquals} from "./utils-albumIdEquals";

export interface SelectAlbumQuery {
    loaded: boolean
    currentAlbumId: AlbumId | undefined
    albumId: AlbumId
}

export interface HasType {
    type: string
}

export type ActionObserver<T extends HasType> = (action: T) => void

export interface MediaPerDayLoaderInterface {
    loadMedias(albumId: AlbumId): Promise<MediasLoadedAction>
}

export class SelectAlbumHandler {

    constructor(
        private readonly mediaPerDayLoader: MediaPerDayLoaderInterface,
    ) {
    }

    public async onSelectAlbum(query: SelectAlbumQuery, dispatch: ActionObserver<MediasLoadedAction | StartLoadingMediasAction>) {
        if (query.loaded && !albumIdEquals(query.currentAlbumId, query.albumId)) {
            dispatch(startLoadingMediasAction(query.albumId))

            return this.mediaPerDayLoader.loadMedias(query.albumId).then(dispatch)
        }
    }
}

