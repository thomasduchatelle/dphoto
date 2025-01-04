import {HasType, SelectAlbumHandler} from "./SelectAlbumHandler";
import {mediasLoadedAction, MediasLoadedAction, startLoadingMediasAction} from "./catalog-actions";
import {MediaPerDayLoader} from "./MediaPerDayLoader";
import {AlbumId, Media, MediaId, MediaType} from "./catalog-state";
import {FetchAlbumMediasPort} from "./CatalogViewerLoader";

function newMedia(mediaId: MediaId, dateTime: string): Media {
    return {
        id: mediaId,
        type: MediaType.IMAGE,
        time: new Date(dateTime),
        uiRelativePath: `${mediaId}/image-${mediaId}.jpg`,
        contentPath: `/content/$\{id}/image-${mediaId}.jpg`,
        source: 'Samsung Galaxy S24'
    };
}

describe('SelectAlbumHandler', () => {
    const albumId = {owner: 'tony@stark.com', folderName: 'avenger-1'};

    let actionObserverFake: ActionObserverFake
    let mediaRepositoryFake: MediaRepositoryFake
    let selectAlbumHandler: SelectAlbumHandler

    beforeEach(() => {
        actionObserverFake = new ActionObserverFake()
        mediaRepositoryFake = new MediaRepositoryFake()
        selectAlbumHandler = new SelectAlbumHandler(new MediaPerDayLoader(mediaRepositoryFake))
    })

    it('it should not start to load if the initial loading is not complete', () => {
        selectAlbumHandler.onSelectAlbum({loaded: false, currentAlbumId: undefined, albumId}, actionObserverFake.onAction);

        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it('it should not start to load if the previous album is same as the current one', () => {
        selectAlbumHandler.onSelectAlbum({loaded: true, currentAlbumId: albumId, albumId}, actionObserverFake.onAction);

        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it('should load medias (grouped by days) and publish 2 actions: MediasLoadedAction and AlbumsAndMediasLoadedAction', async () => {
        const media1 = newMedia('01', "2024-12-01T15:22:00Z");
        const media2 = newMedia('02', "2024-12-01T13:09:00Z");
        const media3 = newMedia('03', "2024-12-02T09:45:00Z");

        mediaRepositoryFake.addMedias(albumId, [
            media1,
            media2,
            media3,
        ] as Media[])

        await selectAlbumHandler.onSelectAlbum({loaded: true, currentAlbumId: undefined, albumId}, actionObserverFake.onAction);

        expect(actionObserverFake.actions).toEqual([
            startLoadingMediasAction(albumId),
            mediasLoadedAction(albumId, [
                {
                    day: new Date("2024-12-01T00:00:00Z"),
                    medias: [media1, media2],
                },
                {
                    day: new Date("2024-12-02T00:00:00Z"),
                    medias: [media3],
                },
            ]),
        ])
    })

    it('should return an empty media list if no medias are found', async () => {
        await selectAlbumHandler.onSelectAlbum({loaded: true, currentAlbumId: undefined, albumId}, actionObserverFake.onAction);

        expect(actionObserverFake.actions).toEqual([
            startLoadingMediasAction(albumId),
            mediasLoadedAction(albumId, []),
        ])
    })

    it('should let error throw', () => {
        const error = new Error("TEST simulate error on loadMedias");
        const selectAlbumHandler = new SelectAlbumHandler({
            loadMedias: (albumId: AlbumId): Promise<MediasLoadedAction> => {
                return Promise.reject(error)
            },
        })

        expect(selectAlbumHandler.onSelectAlbum({loaded: true, currentAlbumId: undefined, albumId}, actionObserverFake.onAction)).rejects.toThrow(error)
        expect(actionObserverFake.actions).toEqual([
            startLoadingMediasAction(albumId),
        ])
    })
});

export class ActionObserverFake {
    public actions: HasType[] = []

    onAction = (action: HasType): void => {
        this.actions.push(action)
    }
}

class MediaRepositoryFake implements FetchAlbumMediasPort {
    private medias: Map<AlbumId, Media[]> = new Map()

    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.resolve(this.medias.get(albumId) ?? [])
    }

    addMedias(albumId: AlbumId, medias: Media[]) {
        this.medias.set(albumId, [...(this.medias.get(albumId) ?? []), ...medias])
    }
}