import {AlbumId, albumIdEquals, Media, MediaWithinADay} from "./catalog-model";
import {mediasLoadedAction, MediasLoadedAction, StartLoadingMediasAction, startLoadingMediasAction} from "./catalog-actions";
import {FetchAlbumMediasPort} from "./CatalogViewerLoader";

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

    public async onSelectAlbum(query: SelectAlbumQuery, observer: ActionObserver<MediasLoadedAction | StartLoadingMediasAction>) {
        if (query.loaded && !albumIdEquals(query.currentAlbumId, query.albumId)) {
            observer(startLoadingMediasAction(query.albumId))

            return this.mediaPerDayLoader.loadMedias(query.albumId).then(observer)
        }
    }
}

export class MediaPerDayLoader {

    constructor(
        private readonly fetchAlbumMediasPort: FetchAlbumMediasPort,
    ) {
    }

    public async loadMedias(albumId: AlbumId): Promise<MediasLoadedAction> {
        const medias = await this.fetchAlbumMediasPort.fetchMedias(albumId)
        return mediasLoadedAction(albumId, this.groupByDay(medias))
    }

    groupByDay(medias: Media[]): MediaWithinADay[] {
        let result: MediaWithinADay[] = []

        medias.forEach(m => {
            const beginning = new Date(m.time)
            beginning.setHours(0, 0, 0, 0)

            if (result.length > 0 && result[0].day.getTime() === beginning.getTime()) {
                result[0].medias.push(m)
            } else {
                result = [{
                    day: beginning,
                    medias: [m],
                }, ...result]
            }
        })

        result.reverse()
        return result
    }
}