import {CreateAlbumHandler, CreateAlbumRequest} from "./CreateAlbumHandler";
import {ActionObserverFake, MediaRepositoryFake, newMedia} from "./SelectAlbumHandler.test";
import {AlbumsAndMediasLoadedAction} from "./catalog-actions";
import {AlbumId, Owner} from "./catalog-state";

const owner1: Owner = "owner1"

describe("CreateAlbumHandler", () => {
    const folderName = "/jan-2025"
    const media1 = newMedia('01', "2024-12-01T15:22:00Z");


    const janAlbumId = {owner: owner1, folderName: folderName}
    const janAlbum = {
        albumId: janAlbumId,
        name: "Jan 2025",
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 1),
        totalCount: 0,
        temperature: 0,
        relativeTemperature: 0,
        sharedWith: []
    }

    let actionObserverFake: ActionObserverFake
    let albumCatalogFake: AlbumCatalogFake
    let mediaRepositoryFake: MediaRepositoryFake
    let handler: CreateAlbumHandler

    beforeEach(() => {
        actionObserverFake = new ActionObserverFake()
        handler = new CreateAlbumHandler()
    })

    it("should raise the AlbumCreatingAction, then a AlbumsAndMediasLoadedAction where the newly created album is selected", async () => {
        mediaRepositoryFake.addMedias(janAlbumId, [media1])

        const createJanAlbumRequest = {
            start: new Date(2025, 1, 1),
            end: new Date(2025, 2, 1),
            forcedFolderName: "",
            name: "Jan 2025"
        }

        await handler.handleCreateAlbum(createJanAlbumRequest)

        expect(albumCatalogFake.albumCreationRequests).toEqual([
            {
                start: new Date(2025, 1, 1),
                end: new Date(2025, 2, 1),
                forcedFolderName: "",
                name: "Jan 2025"
            }
        ])

        expect(actionObserverFake.actions).toEqual([
            {
                type: "AlbumCreatingAction",
            },
            {
                type: "AlbumsAndMediasLoadedAction",
                albums: [
                    janAlbum
                ],
                medias: [
                    {
                        day: new Date("2024-12-01T00:00:00Z"),
                        medias: [media1],
                    },
                ],
                selectedAlbum: janAlbum
            } as AlbumsAndMediasLoadedAction,
        ])
    })

    // it should return an error if the rest adapter fails to process the request with the error code from the request
})

class AlbumCatalogFake {
    albumCreationRequests: CreateAlbumRequest[] = []

    createAlbum(request: CreateAlbumRequest): Promise<AlbumId> {
        this.albumCreationRequests.push(request) // TODO where is the owner coming from ?
        return Promise.resolve({owner: owner1, folderName: `/${request.name.replaceAll(" ", "-").toLowerCase()}`})
    }
}

export {}