import {HasType} from "./ActionObserver";
import {Album, AlbumId, Media, MediaId, MediaType} from "./catalog-state";
import {CatalogLoader, FetchAlbumsPort, PartialCatalogLoaderState} from "./CatalogLoader";
import {FetchAlbumMediasPort, MediaPerDayLoader} from "./MediaPerDayLoader";
import {AlbumsAndMediasLoadedAction, MediasLoadedAction} from "./catalog-actions";


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

const expectedMediasPerDay = [
    {
        day: new Date(2024, 11, 1),
        medias: [medias[0], medias[1]],
    },
    {
        day: new Date(2024, 11, 2),
        medias: [medias[2]],
    },
];

function repositoryWithThreeMedias() {
    const map = new Map<AlbumId, Media[]>();

    map.set(twoAlbums[0].albumId, medias)
    return new MediaPerDayLoader(new MediaRepositoryFake(map));

}

describe("CatalogLoader", () => {
    const stateNotLoaded: PartialCatalogLoaderState = {albumsLoaded: false, allAlbums: []}
    const firstAlbumLoaded: PartialCatalogLoaderState = {albumsLoaded: true, allAlbums: twoAlbums, mediasLoadedFromAlbumId: twoAlbums[0].albumId}
    const firstAlbumLoading: PartialCatalogLoaderState = {albumsLoaded: false, allAlbums: twoAlbums, loadingMediasFor: twoAlbums[0].albumId}

    const newCatalogLoader = (dispatch: ActionObserverFake, state: PartialCatalogLoaderState) => {
        return new CatalogLoader(dispatch.onAction, repositoryWithThreeMedias(), new AlbumRepositoryFake(twoAlbums), state)
    }

    const newLoaderFailingGettingMedias = (dispatch: ActionObserverFake, error: Error, state: PartialCatalogLoaderState) => {
        return new CatalogLoader(dispatch.onAction, new MediaPerDayLoader(new MediaRepositoryFakeWithFailure(error)), new AlbumRepositoryFake(twoAlbums), state)
    }

    let dispatch: ActionObserverFake

    beforeEach(() => {
        dispatch = new ActionObserverFake()
    })

    it("should load albums and medias of requested album when the state is not loaded", async () => {
        const loader = newCatalogLoader(dispatch, stateNotLoaded)
        await loader.onPageRefresh(twoAlbums[0].albumId)

        expect(dispatch.actions).toEqual([
            {
                type: 'AlbumsAndMediasLoadedAction',
                albums: twoAlbums,
                medias: expectedMediasPerDay,
                selectedAlbum: twoAlbums[0],
            } as AlbumsAndMediasLoadedAction,
        ])
    })

    it("should dispatch a MediaFailedToLoadAction if the medias cannot be loaded", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error, stateNotLoaded)
        await loader.onPageRefresh(twoAlbums[0].albumId)

        expect(dispatch.actions).toEqual([
            {
                type: 'MediaFailedToLoadAction',
                albums: twoAlbums,
                selectedAlbum: twoAlbums[0],
                error: new Error(`failed to load medias of ${twoAlbums[0].albumId}`, error),
            }
        ])
    })

    it("should throw an error if the albums cannot be loaded", async () => {
        const error = new Error("TEST simulate albums failing to load");
        const loader = new CatalogLoader(dispatch.onAction, new MediaPerDayLoader(new MediaRepositoryFake()), new AlbumRepositoryFakeWithError(error), stateNotLoaded)
        await expect(loader.onPageRefresh(twoAlbums[0].albumId)).rejects.toThrow(error)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should not load if the state is currently loading the desired album", async () => {
        const loader = newCatalogLoader(dispatch, firstAlbumLoading)
        await loader.onPageRefresh(twoAlbums[0].albumId)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should load albums and medias of the first album when the state is not loaded and there is no album selected", async () => {
        const loader = newCatalogLoader(dispatch, stateNotLoaded)
        await loader.onPageRefresh(undefined)

        expect(dispatch.actions).toEqual([
            {
                type: 'AlbumsAndMediasLoadedAction',
                albums: twoAlbums,
                medias: expectedMediasPerDay,
                selectedAlbum: twoAlbums[0],
                redirectTo: twoAlbums[0].albumId,
            } as AlbumsAndMediasLoadedAction,
        ])
    })

    it("should not trigger an 'initial load' if the albums are already loaded", async () => {
        const loader = newCatalogLoader(dispatch, firstAlbumLoaded)
        await loader.onPageRefresh(undefined)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should dispatch a MediaFailedToLoadAction when failing to load the medias of the selected album", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error, stateNotLoaded)
        await loader.onPageRefresh(undefined)

        expect(dispatch.actions).toEqual([
            {
                type: 'MediaFailedToLoadAction',
                albums: twoAlbums,
                selectedAlbum: twoAlbums[0],
                error: new Error(`failed to load medias of ${twoAlbums[0].albumId}`, error),
            },
        ])
    })

    it("should not do anything if the medias are already loaded for the requested album", async () => {
        const loader = newCatalogLoader(dispatch, firstAlbumLoaded)
        await loader.onPageRefresh(twoAlbums[0].albumId)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should load the medias if another album is selected", async () => {
        const loader = newCatalogLoader(dispatch, firstAlbumLoaded)
        await loader.onPageRefresh(twoAlbums[1].albumId)

        expect(dispatch.actions).toEqual([
            {
                type: 'MediasLoadedAction',
                albumId: twoAlbums[1].albumId,
                medias: [],
            } as MediasLoadedAction,
        ])
    })

    it("should not load the medias if another album is selected but the medias are already loading", async () => {
        const loader = newCatalogLoader(dispatch, {
            ...firstAlbumLoaded,
            loadingMediasFor: twoAlbums[1].albumId,
        })
        await loader.onPageRefresh(twoAlbums[1].albumId)

        expect(dispatch.actions).toHaveLength(0)
    })

    it("should dispatch an error if the medias cannot be loaded", async () => {
        const error = new Error("TEST simulate medias failing to load");
        const loader = newLoaderFailingGettingMedias(dispatch, error, firstAlbumLoaded)
        await loader.onPageRefresh(twoAlbums[1].albumId)

        expect(dispatch.actions).toEqual([
            {
                type: 'MediaFailedToLoadAction',
                selectedAlbum: twoAlbums[1],
                error,
            },
        ])
    })
})

class ActionObserverFake {
    public actions: HasType[] = []

    onAction = (action: HasType): void => {
        this.actions.push(action)
    }
}

class AlbumRepositoryFake implements FetchAlbumsPort {
    constructor(
        private albums: Album[] = []
    ) {
    }

    fetchAlbums = (): Promise<Album[]> => {
        return Promise.resolve(this.albums)
    }
}

class AlbumRepositoryFakeWithError implements FetchAlbumsPort {
    constructor(
        private error: Error,
    ) {
    }

    fetchAlbums = (): Promise<Album[]> => {
        return Promise.reject(this.error)
    }
}

class MediaRepositoryFake implements FetchAlbumMediasPort {
    constructor(
        private medias: Map<AlbumId, Media[]> = new Map()
    ) {
    }


    fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return Promise.resolve(this.medias.get(albumId) ?? [])
    }
}

class MediaRepositoryFakeWithFailure implements FetchAlbumMediasPort {
    constructor(
        private error: Error,
    ) {
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