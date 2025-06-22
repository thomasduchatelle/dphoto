import {DPhotoApplication} from "../application";
import {CatalogAPIAdapter} from "./adapters/api";
import {CatalogFactory} from "./catalog-factories";
import {rest} from "msw";
import {SetupServer, setupServer} from "msw/node";
import {MediaType} from "./language";
import {CatalogViewerAction} from "./actions";
import {albumsAndMediasLoaded, MediaPerDayLoader, MediasLoaded, OnPageRefresh} from "./navigation";

describe('CatalogFactory', () => {

    const albumIdAvenger1 = {
        folderName: "avenger-1",
        owner: "tony@stark.com"
    }

    let server: SetupServer

    beforeAll(() => {
        server = setupServer()
        server.listen()
    })

    beforeEach(() => {
        server.resetHandlers()
    })

    afterAll(() => {
        server.close()
    })

    it('should create a new instance of CatalogAPIAdapter', () => {
        const restAdapter = newCatalogFactory().restAdapter();
        expect(restAdapter).toBeInstanceOf(CatalogAPIAdapter);
    });

    it('should create a new instance of CatalogViewerLoader', async () => {
        const dispatch: CatalogViewerAction[] = []
        const restAdapter = newCatalogFactory().restAdapter();
        const mediaViewLoader = new OnPageRefresh(dispatch.push.bind(dispatch), new MediaPerDayLoader(restAdapter), restAdapter);
        expect(mediaViewLoader).toBeInstanceOf(OnPageRefresh);

        server.use(
            getAlbums(avenger1Album()),
            getOwners(ownerTonyStark()),
            getMediasForAvenger1(),
        )

        await mediaViewLoader.onPageRefresh({
            mediasLoadedFromAlbumId: undefined,
            allAlbums: [],
            albumsLoaded: false,
            loadingMediasFor: undefined,
        });
        expect(dispatch).toEqual([
            albumsAndMediasLoaded({
                albums: [
                    {
                        albumId: albumIdAvenger1,
                        end: new Date(2021, 0, 31),
                        name: "Avenger 1",
                        ownedBy: {
                            name: "Tony Stark",
                            users: [
                                {
                                    email: "tony@stark.com",
                                    name: "Tony Stark",
                                    picture: "http://tony-stark.com/picture.jpg"
                                }
                            ]
                        },
                        relativeTemperature: 1,
                        sharedWith: [],
                        start: new Date(2021, 0, 1),
                        temperature: 0.3333333333333333,
                        totalCount: 10
                    }
                ],
                selectedAlbum: {
                    albumId: albumIdAvenger1,
                    end: new Date(2021, 0, 31),
                    name: "Avenger 1",
                    ownedBy: {
                        name: "Tony Stark",
                        users: [
                            {
                                email: "tony@stark.com",
                                name: "Tony Stark",
                                picture: "http://tony-stark.com/picture.jpg"
                            }
                        ]
                    },
                    relativeTemperature: 1,
                    sharedWith: [],
                    start: new Date(2021, 0, 1),
                    temperature: 0.3333333333333333,
                    totalCount: 10
                },
                medias: [
                    {
                        day: new Date(2021, 0, 5),
                        medias: [
                            {
                                contentPath: "/api/v1/owners/tony@stark.com/medias/media-1/image.jpg?access_token=",
                                id: "media-1",
                                source: "Ironman Suit",
                                time: new Date("2021-01-05T12:42:00Z"),
                                type: MediaType.IMAGE,
                                uiRelativePath: "media-1/image.jpg"
                            }
                        ]
                    }
                ],
                redirectTo: albumIdAvenger1,
            })]);
    });

    it('should create a new instance of SelectAlbumHandler', async () => {
        const dispatch: CatalogViewerAction[] = []
        const restAdapter = newCatalogFactory().restAdapter();
        const mediaViewLoader = new OnPageRefresh(dispatch.push.bind(dispatch), new MediaPerDayLoader(restAdapter), restAdapter);

        server.use(
            getMediasForAvenger1(),
        )

        await mediaViewLoader.onPageRefresh({
            mediasLoadedFromAlbumId: undefined,
            allAlbums: [],
            albumsLoaded: true,
            loadingMediasFor: undefined,
            albumId: albumIdAvenger1,
        });

        let mediasLoadedAction = dispatch.find(action => action.type === "mediasLoaded");
        expect(mediasLoadedAction).toBeDefined()

        expect((mediasLoadedAction as MediasLoaded).payload?.medias).toEqual(
            [
                {
                    day: new Date(2021, 0, 5),
                    medias: [
                        {
                            contentPath: "/api/v1/owners/tony@stark.com/medias/media-1/image.jpg?access_token=",
                            id: "media-1",
                            source: "Ironman Suit",
                            time: new Date("2021-01-05T12:42:00Z"),
                            type: MediaType.IMAGE,
                            uiRelativePath: "media-1/image.jpg"
                        }
                    ]
                }
            ]
        );
    });

    function newCatalogFactory() {
        return new CatalogFactory(newDPhotoApplication());
    }

    function newDPhotoApplication() {
        return new DPhotoApplication();
    }
});

function avenger1Album() {
    return {
        owner: "tony@stark.com",
        folderName: "/avenger-1",
        name: "Avenger 1",
        start: "2021-01-01",
        end: "2021-01-31",
        totalCount: 10,
    };
}

function getAlbums(...albums: any[]) {
    return rest.get('/api/v1/albums', (req, res, ctx) => {

        return res(
            ctx.json(albums)
        )
    });
}

function ownerTonyStark() {
    return {
        id: "tony@stark.com",
        name: "Tony Stark",
        users: [
            {
                name: "Tony Stark",
                email: "tony@stark.com",
                picture: "http://tony-stark.com/picture.jpg",
            }
        ],
    };
}

function getOwners(...owners: any[]) {
    return rest.get('/api/v1/owners', (req, res, ctx) => {

        return res(
            ctx.json(owners)
        )
    });
}

function getMediasForAvenger1() {
    return rest.get('/api/v1/owners/tony@stark.com/albums/avenger-1/medias', (req, res, ctx) => {

        return res(
            ctx.json([
                {
                    id: "media-1",
                    type: "IMAGE",
                    filename: "image.jpg",
                    time: "2021-01-05T12:42:00Z",
                    source: "Ironman Suit",
                },
            ])
        )
    });
}