import {ActionObserverFake, MediaRepositoryFake, newMedia} from "./SelectAlbumHandler.test";
import {AlbumId, Owner} from "./catalog-state";
import {
    albumFolderNameAlreadyTakenErr,
    albumStartAndEndDateMandatoryErr,
    CreateAlbumController,
    CreateAlbumRequest,
    CreateAlbumState,
    emptyCreateAlbum
} from "./CreateAlbumController";
import dayjs from "dayjs";
import {CatalogError} from "./errors";

const owner1: Owner = "owner1"

describe("CreateAlbumController", () => {
    const defaultDateForEmptyAlbum = dayjs("2025-01-04");
    const folderName = "/avenger_3"

    const stateValidForSubmission = {
        name: "Avenger 3",
        start: dayjs("2025-01-04T09:08:07"),
        startsAtStartOfTheDay: true,
        end: dayjs("2025-01-12T13:42:12"),
        endsAtEndOfTheDay: true,
        forceFolderName: "",
        withCustomFolderName: false,
        open: true
    }

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

    let stateHolder: StateHolder<CreateAlbumState>
    let handler: CreateAlbumController

    beforeEach(() => {
        stateHolder = new StateHolder(emptyCreateAlbum(defaultDateForEmptyAlbum))
        actionObserverFake = new ActionObserverFake()
        albumCatalogFake = new AlbumCatalogFake()

        handler = new CreateAlbumController(
            stateHolder.update,
            albumCatalogFake,
            actionObserverFake.onAction,
            defaultDateForEmptyAlbum,
        )
    })

    it("it should reset the state when opening the dialog to create a new album", () => {
        stateHolder.state = {
            open: false,
            name: "Jan 2025",
            start: dayjs("2025-01-01"),
            end: dayjs("2025-02-01"),
            forceFolderName: "",
            startsAtStartOfTheDay: true,
            endsAtEndOfTheDay: true,
            withCustomFolderName: false,
        }

        handler.openDialogForCreateAlbum()

        expect(stateHolder.state).toEqual({
            ...emptyCreateAlbum(defaultDateForEmptyAlbum),
            open: true,
        })
        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it("it should close the dialog without any other actions when onClose is requested", () => {
        stateHolder.state.open = true

        handler.onCloseCreateAlbumDialog()

        expect(stateHolder.state.open).toEqual(false)
        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it.each([
        [null, dayjs()],
        [dayjs(), null],
        [dayjs(), dayjs().subtract(1, "day")],
    ])("it should not allow to submit the form if the dates are invalid and display albumStartAndEndDateMandatoryErr error instead: %s -> %s", async (start, end) => {
        await handler.onSubmitCreateAlbum({
            ...stateHolder.state,
            start,
            end,
        })

        expect(albumCatalogFake.albumCreationRequests).toHaveLength(0)
        expect(stateHolder.state.errorCode).toEqual(albumStartAndEndDateMandatoryErr)
        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it("it should submit the form with requested start date and exclusive end date, then dispatch an action with albums and medias reloaded and close the dialog.", async () => {
        await handler.onSubmitCreateAlbum({
            name: "Avenger 3",
            start: dayjs("2025-01-04T09:08:07"),
            startsAtStartOfTheDay: true,
            end: dayjs("2025-01-12T13:42:12"),
            endsAtEndOfTheDay: true,
            forceFolderName: "/marvel_2018_avenger_3",
            withCustomFolderName: false,
            open: true
        })

        expect(albumCatalogFake.albumCreationRequests).toEqual([
            {
                name: "Avenger 3",
                start: new Date(2025, 0, 4),
                end: new Date(2025, 0, 13),
                forcedFolderName: "",
            }
        ])

        // expect(stateHolder.state.creating).toEqual(true) // TODO adds the state
        expect(stateHolder.state.open).toEqual(false)
        expect(actionObserverFake.actions).toHaveLength(0) // TODO it should have something - eventually
    })

    it("it should submit the form with custom start date and end date", async () => {
        await handler.onSubmitCreateAlbum({
            name: "Avenger 3",
            start: dayjs("2025-01-04T09:08:07"),
            startsAtStartOfTheDay: false,
            end: dayjs("2025-01-12T13:42:12"),
            endsAtEndOfTheDay: false,
            forceFolderName: "",
            withCustomFolderName: false,
            open: true
        })

        expect(albumCatalogFake.albumCreationRequests).toEqual([
            {
                name: "Avenger 3",
                start: new Date(2025, 0, 4, 9, 8, 7),
                end: new Date(2025, 0, 12, 13, 42, 12),
                forcedFolderName: "",
            }
        ])
    })

    it("it should submit the form with forced folderName", async () => {
        await handler.onSubmitCreateAlbum({
            name: "Avenger 3",
            start: dayjs("2025-01-04T09:08:07"),
            startsAtStartOfTheDay: true,
            end: dayjs("2025-01-12T13:42:12"),
            endsAtEndOfTheDay: true,
            forceFolderName: "/marvel_2018_avenger_3",
            withCustomFolderName: true,
            open: true
        })

        expect(albumCatalogFake.albumCreationRequests).toEqual([
            {
                name: "Avenger 3",
                start: new Date(2025, 0, 4),
                end: new Date(2025, 0, 13),
                forcedFolderName: "/marvel_2018_avenger_3",
            }
        ])

        expect(stateHolder.state.open).toEqual(false)
    })

    it("should display the business error when saving failed", async () => {
        albumCatalogFake.failsWithError = {
            errorCode: albumFolderNameAlreadyTakenErr,
            message: "TEST error",
        } as CatalogError

        stateHolder.state = stateValidForSubmission
        await handler.onSubmitCreateAlbum(stateHolder.state)

        expect(stateHolder.state).toEqual({
            ...stateValidForSubmission,
            errorCode: albumFolderNameAlreadyTakenErr,
            open: true,
            // loading: false, // TODO
        })
        expect(actionObserverFake.actions).toHaveLength(0)
    })

    it("should fallback on a native error when saving fails for an unknown reason", async () => {
        albumCatalogFake.failsWithError = new Error("TEST Create Album Error")

        stateHolder.state = stateValidForSubmission
        await handler.onSubmitCreateAlbum(stateHolder.state)

        expect(stateHolder.state).toEqual({
            ...stateValidForSubmission,
            errorCode: "TEST Create Album Error",
            open: true,
            // loading: false, // TODO
        })
        expect(actionObserverFake.actions).toHaveLength(0)
    })

    // it should switch to the newly created album after the creation

    // it("should raise the AlbumCreatingAction, then a AlbumsAndMediasLoadedAction where the newly created album is selected", async () => {
    //     mediaRepositoryFake.addMedias(janAlbumId, [media1])
    //
    //     const createJanAlbumRequest = {
    //         start: new Date(2025, 1, 1),
    //         end: new Date(2025, 2, 1),
    //         forcedFolderName: "",
    //         name: "Jan 2025"
    //     }
    //
    //     await handler.handleCreateAlbum(createJanAlbumRequest)
    //
    //     expect(albumCatalogFake.albumCreationRequests).toEqual([
    //         {
    //             start: new Date(2025, 1, 1),
    //             end: new Date(2025, 2, 1),
    //             forcedFolderName: "",
    //             name: "Jan 2025"
    //         }
    //     ])
    //
    //     expect(actionObserverFake.actions).toEqual([
    //         {
    //             type: "AlbumCreatingAction",
    //         },
    //         {
    //             type: "AlbumsAndMediasLoadedAction",
    //             albums: [
    //                 janAlbum
    //             ],
    //             medias: [
    //                 {
    //                     day: new Date("2024-12-01T00:00:00Z"),
    //                     medias: [media1],
    //                 },
    //             ],
    //             selectedAlbum: janAlbum
    //         } as AlbumsAndMediasLoadedAction,
    //     ])
    // })

    // it should return an error if the rest adapter fails to process the request with the error code from the request
})

class StateHolder<T> {
    constructor(
        public state: T,
    ) {
    }

    update = (stateUpdater: (prev: T) => T) => {
        this.state = stateUpdater(this.state)
    }
}

class AlbumCatalogFake {
    albumCreationRequests: CreateAlbumRequest[] = []
    failsWithError: Error | undefined

    createAlbum(request: CreateAlbumRequest): Promise<AlbumId> {
        if (this.failsWithError) {
            return Promise.reject(this.failsWithError)
        }

        this.albumCreationRequests.push(request); // TODO where is the owner coming from ?
        return Promise.resolve({owner: owner1, folderName: `/${request.name.replaceAll(" ", "-").toLowerCase()}`})
    }
}

export {}