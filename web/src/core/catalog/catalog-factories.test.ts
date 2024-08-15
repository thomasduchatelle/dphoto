import {DPhotoApplication} from "../application";
import {AlbumsAndMediasLoadedAction, CatalogViewerLoader, MediasLoadedAction, MediaType} from "./index";
import {CatalogAPIAdapter} from "./adapters/api";
import {CatalogFactory} from "./catalog-factories";
import {rest} from "msw";
import {SetupServer, setupServer} from "msw/node";
import {ActionObserverFake} from "./domain/SelectAlbumHandler.test";

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
        const mediaViewLoader = newCatalogFactory().mediaViewLoader();
        expect(mediaViewLoader).toBeInstanceOf(CatalogViewerLoader);

        server.use(
            getAlbums(avenger1Album()),
            getOwners(ownerTonyStark()),
            getMediasForAvenger1(),
        )

        const got = await mediaViewLoader.loadInitialCatalog({});
        expect(got).toEqual({
            type: 'AlbumsAndMediasLoadedAction',
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
            media: [
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
        } as AlbumsAndMediasLoadedAction);
    });

    it('should create a new instance of SelectAlbumHandler', async () => {
        const selectAlbumHandler = newCatalogFactory().selectAlbumHandler();

        server.use(
            getMediasForAvenger1(),
        )

        const observer = new ActionObserverFake();
        await selectAlbumHandler.onSelectAlbum({
            loaded: true,
            albumId: albumIdAvenger1,
            currentAlbumId: undefined,
        }, observer.onAction);

        let mediasLoadedAction = observer.actions.find(action => action.type === "MediasLoadedAction");
        expect(mediasLoadedAction).toBeDefined()

        expect((mediasLoadedAction as MediasLoadedAction).medias).toEqual(
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