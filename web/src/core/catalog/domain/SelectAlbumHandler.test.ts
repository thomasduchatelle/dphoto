import {SelectAlbumHandler} from "./SelectAlbumHandler";
import {mediasLoadedAction, MediasLoadedAction, startLoadingMediasAction} from "./catalog-actions";
import {MediaPerDayLoader} from "./MediaPerDayLoader";
import {AlbumId, Media, MediaId, MediaType} from "./catalog-state";
import {FetchAlbumMediasPort} from "./CatalogViewerLoader";
import {HasType} from "./ActionObserver";

export function newMedia(mediaId: MediaId, dateTime: string): Media {
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
    const anotherAlbumId = {owner: 'pepper@stark.com', folderName: '/stark-industries'};

    const media1 = newMedia('01', "2024-12-01T15:22:00Z");
    const media2 = newMedia('02', "2024-12-01T13:09:00Z");
    const media3 = newMedia('03', "2024-12-02T09:45:00Z");

    const addsMediasToMediaRepositoryFake = () => mediaRepositoryFake.addMedias(albumId, [
        media1,
        media2,
        media3,
    ] as Media[])

    let actionObserverFake: ActionObserverFake
    let mediaRepositoryFake: MediaRepositoryFake

    beforeEach(() => {
        actionObserverFake = new ActionObserverFake()
        mediaRepositoryFake = new MediaRepositoryFake()
    })

    it('it should not start to load if the initial loading is not complete', () => {
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake))
        selectAlbumHandler.onSelectAlbum(albumId);

        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it('it should not start to load if the previous album is same as the current one', () => {
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake), albumId)
        selectAlbumHandler.onSelectAlbum(albumId);

        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it('should load medias (grouped by days) and publish 2 actions: MediasLoadedAction and AlbumsAndMediasLoadedAction', async () => {
        addsMediasToMediaRepositoryFake()
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake), anotherAlbumId)
        await selectAlbumHandler.onSelectAlbum(albumId);

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

    it('should load medias if the loading album is different from the one requested (even if the one currently selected is the target)', async () => {
        addsMediasToMediaRepositoryFake()
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake), albumId, anotherAlbumId)
        await selectAlbumHandler.onSelectAlbum(albumId);

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

    it('should not load medias if there is no album currently loaded (to prevent double trigger of the REST call)', async () => {
        addsMediasToMediaRepositoryFake()
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake), undefined, anotherAlbumId)
        await selectAlbumHandler.onSelectAlbum(albumId);

        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it('should return an empty media list if no medias are found', async () => {
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, new MediaPerDayLoader(mediaRepositoryFake), anotherAlbumId)
        await selectAlbumHandler.onSelectAlbum(albumId);

        expect(actionObserverFake.actions).toEqual([
            startLoadingMediasAction(albumId),
            mediasLoadedAction(albumId, []),
        ])
    })

    it('should let error throw', () => {
        const error = new Error("TEST simulate error on loadMedias");
        const loaderThrowingAnError = {
            loadMedias: (albumId: AlbumId): Promise<MediasLoadedAction> => {
                return Promise.reject(error)
            },
        }
        const selectAlbumHandler = new SelectAlbumHandler(actionObserverFake.onAction, loaderThrowingAnError, anotherAlbumId)

        expect(selectAlbumHandler.onSelectAlbum(albumId)).rejects.toThrow(error)
        expect(actionObserverFake.actions).toEqual([
            startLoadingMediasAction(albumId),
        ])
    })
})

export class ActionObserverFake {
    public actions: HasType[] = []

    onAction = (action: HasType): void => {
        this.actions.push(action)
    }
}

export class MediaRepositoryFake implements FetchAlbumMediasPort {
    private medias: Map<AlbumId, Media[]> = new Map()

    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.resolve(this.medias.get(albumId) ?? [])
    }

    addMedias(albumId: AlbumId, medias: Media[]) {
        this.medias.set(albumId, [...(this.medias.get(albumId) ?? []), ...medias])
    }
}