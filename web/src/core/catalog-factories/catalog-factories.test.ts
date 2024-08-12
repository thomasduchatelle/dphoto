import {DPhotoApplication} from "../application";
import {Album, AlbumsAndMediasLoadedAction, MediaViewLoader, MediaWithinADay} from "../catalog";
import {CatalogAPIAdapter} from "../catalog-apis";
import {CatalogFactory} from "./catalog-factories";
import {http, HttpResponse} from "msw";
import {setupServer} from "msw/node";

// write tests with jest for the CatalogFactory class

let application: DPhotoApplication
let catalogFactory: CatalogFactory
describe('CatalogFactory', () => {

    beforeEach(() => {
        application = new DPhotoApplication();
        catalogFactory = new CatalogFactory(application);
    });

    it('should create a new instance of CatalogAPIAdapter', () => {
        const restAdapter = catalogFactory.restAdapter();
        expect(restAdapter).toBeInstanceOf(CatalogAPIAdapter);
    });

    it('should create a new instance of MediaViewLoader', async () => {
        const mediaViewLoader = newCatalogFactory().mediaViewLoader();
        expect(mediaViewLoader).toBeInstanceOf(MediaViewLoader);

        const handlers = [
            http.get('https://example.com/user', () => {
                // ...and respond to them using this JSON response.
                return HttpResponse.json({
                    id: 'c7b3d8e0-5e0b-4b0f-8b3a-3b9f4b3d3b3d',
                    firstName: 'John',
                    lastName: 'Maverick',
                })
            }),
        ]
        const server = setupServer(...handlers)
        // server.listen()
        // server.use(http.get('https://example.com/user', (req, res, ctx) => {
        //     return res(ctx.json({id: 'c7b3d8e0-5e0b-4b0f-8b3a-3b9f4b3d3b3d', firstName: 'John', lastName: 'Maverick'}))
        // }))

        const got = await mediaViewLoader.loadInitialCatalog({});
        // expect(got).toEqual({
        //     type: 'AlbumsAndMediasLoadedAction',
        //     albums: [],
        //     media: [],
        //     selectedAlbum: undefined,
        // } as AlbumsAndMediasLoadedAction);

        // server.close()

    });

    function newCatalogFactory() {
        return new CatalogFactory(newDPhotoApplication());
    }

    function newDPhotoApplication() {
        return new DPhotoApplication();
    }
});