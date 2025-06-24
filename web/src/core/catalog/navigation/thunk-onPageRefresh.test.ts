import {Album, AlbumId, Media, MediaId, MediaType} from "../language";
import {mediasLoaded} from "./action-mediasLoaded";
import {albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {mediaLoadFailed} from "./action-mediaLoadFailed";
import {FetchAlbumsAndMediasPort, OnPageRefresh, OnPageRefreshArgs} from "./thunk-onPageRefresh";
import {CatalogViewerAction} from "../actions";


const twoAlbums: Album[] = [
    {
        albumId: {owner: "owner1", folderName: "jan-25"},
        name: "January 2025",
        start: new Date(2025, 0, 1),
        end: new Date(2025, 0, 31),
        totalCount: 42,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    },
    {
        albumId: {owner: "owner2", folderName: "feb-25"},
        name: "February 2025",
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 0),
        totalCount: 12,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    },
]

const medias = [
    newMedia('01', "2024-12-01T15:22:00Z"),
    newMedia('02', "2024-12-01T13:09:00Z"),
    newMedia('03', "2024-12-02T09:45:00Z"),
];

interface PartialCatalogLoaderState extends Omit<OnPageRefreshArgs, "albumId"> {
}

describe("CatalogLoader", () => {
    const stateNotLoaded: PartialCatalogLoaderState = {albumsLoaded: false, allAlbums: []}
    const firstAlbumLoaded: PartialCatalogLoaderState = {albumsLoaded: true, allAlbums: twoAlbums, mediasLoadedFromAlbumId: twoAlbums[0].albumId}
    const firstAlbumLoading: PartialCatalogLoaderState = {albumsLoaded: false, allAlbums: twoAlbums, loadingMediasFor: twoAlbums[0].albumId}

    const newThunk = (dispatch: ActionObserverFake) => {
        return new OnPageRefresh(dispatch.onAction, new AlbumAndMediaRepositoryFake(
            twoAlbums,
            new Map([[twoAlbums[0].albumId, medias]])
        ))
    }

    const newLoaderFailingGettingMedias = (dispatch: ActionObserverFake, error: Error) => {
        return new OnPageRefresh(dispatch.onAction, new AlbumAndMediaRepositoryFakeWithMediaFailure(twoAlbums, error))
    }

    let dispatch: ActionObserverFake

    beforeEach(() => {
        dispatch = new ActionObserverFake()
    })

    it("should load albums and medias of requested album when the state is not loaded", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh({...stateNotLoaded, albumId: twoAlbums[0].albumId})

        expect(dispatch.actions).toEqual([
            albumsAndMediasLoaded({
                albums: twoAlbums,
                medias: medias,
                mediasFromAlbumId: twoAlbums[0].albumId,
            }),
        ])
    })

    it("should dispatch a mediaLoadFailedAction if the medias cannot be loaded", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error)
        await loader.onPageRefresh({...stateNotLoaded, albumId: twoAlbums[0].albumId})

        expect(dispatch.actions).toEqual([
            mediaLoadFailed({
                albums: twoAlbums,
                displayedAlbumId: twoAlbums[0].albumId,
                error: new Error(`failed to load medias of ${twoAlbums[0].albumId}`, error),
            })
        ])
    })

    it("should throw an error if the albums cannot be loaded", async () => {
        const error = new Error("TEST simulate albums failing to load");
        const loader = new OnPageRefresh(dispatch.onAction, new AlbumAndMediaRepositoryFakeWithAlbumFailure(error))
        await expect(loader.onPageRefresh({...stateNotLoaded, albumId: twoAlbums[0].albumId})).rejects.toThrow(error)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should not load if the state is currently loading the desired album", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh({...firstAlbumLoading, albumId: twoAlbums[0].albumId})

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should load albums and medias of the first album when the state is not loaded and there is no album selected", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh(stateNotLoaded)

        expect(dispatch.actions).toEqual([
            albumsAndMediasLoaded({
                albums: twoAlbums,
                medias: medias,
                mediasFromAlbumId: twoAlbums[0].albumId,
                redirectTo: twoAlbums[0].albumId,
            }),
        ])
    })

    it("should not trigger an 'initial load' if the albums are already loaded", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh(firstAlbumLoaded)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should dispatch a mediaLoadFailedAction when failing to load the medias of the selected album", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error)
        await loader.onPageRefresh(stateNotLoaded)

        expect(dispatch.actions).toEqual([
            mediaLoadFailed({
                albums: twoAlbums,
                displayedAlbumId: twoAlbums[0].albumId,
                error: new Error(`failed to load medias of ${twoAlbums[0].albumId}`, error),
            }),
        ])
    })

    it("should not do anything if the medias are already loaded for the requested album", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh({...firstAlbumLoaded, albumId: twoAlbums[0].albumId})

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should load the medias if another album is selected", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh({...firstAlbumLoaded, albumId: twoAlbums[1].albumId})

        expect(dispatch.actions).toEqual([
            mediasLoaded({
                albumId: twoAlbums[1].albumId,
                medias: [],
            }),
        ])
    })

    it("should not load the medias if another album is selected but the medias are already loading", async () => {
        const loader = newThunk(dispatch)
        await loader.onPageRefresh({
            ...firstAlbumLoaded,
            loadingMediasFor: twoAlbums[1].albumId,
            albumId: twoAlbums[1].albumId,
        })

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should dispatch an error if the medias cannot be loaded", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error)
        await loader.onPageRefresh({...firstAlbumLoaded, albumId: twoAlbums[1].albumId})

        expect(dispatch.actions).toEqual([
            mediaLoadFailed({
                displayedAlbumId: twoAlbums[1].albumId,
                error,
            }),
        ])
    })
})

class ActionObserverFake {
    public actions: CatalogViewerAction[] = []

    onAction = (action: CatalogViewerAction): void => {
        this.actions.push(action)
    }
}

class AlbumAndMediaRepositoryFake implements FetchAlbumsAndMediasPort {
    constructor(
        private albums: Album[] = [],
        private medias: Map<AlbumId, Media[]> = new Map(),
    ) {
    }

    fetchAlbums = (): Promise<Album[]> => {
        return Promise.resolve(this.albums)
    }

    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.resolve(this.medias.get(albumId) || [])
    }
}

class AlbumAndMediaRepositoryFakeWithAlbumFailure implements FetchAlbumsAndMediasPort {
    constructor(
        private error: Error,
    ) {
    }

    fetchAlbums = (): Promise<Album[]> => {
        return Promise.reject(this.error)
    }

    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.resolve([])
    }
}

class AlbumAndMediaRepositoryFakeWithMediaFailure implements FetchAlbumsAndMediasPort {
    constructor(
        private albums: Album[] = [],
        private error: Error,
    ) {
    }

    fetchAlbums = (): Promise<Album[]> => {
        return Promise.resolve(this.albums)
    }

    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.reject(this.error)
    }
}

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

export {}
