import {UpdateAlbumDatesPort, updateAlbumDatesThunk} from "./thunk-updateAlbumDates";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {albumDatesUpdateFailed} from "./action-albumDatesUpdateFailed";
import {someMedias, twoAlbums} from "../tests/test-helper-state";
import {Album, AlbumId, Media} from "../language";

class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
    public updatedAlbums: { albumId: AlbumId, startDate: Date, endDate: Date }[] = [];
    public shouldFailUpdate = false;
    public shouldFailFetch = false;

    constructor(
        private albums: Album[] = [],
        private medias: Media[] = []
    ) {
    }

    async updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void> {
        if (this.shouldFailUpdate) {
            throw new Error("Failed to update album dates");
        }
        this.updatedAlbums.push({albumId, startDate, endDate});
    }

    async fetchAlbums(): Promise<Album[]> {
        if (this.shouldFailFetch) {
            throw new Error("Failed to fetch albums");
        }
        return this.albums;
    }

    async fetchMedias(albumId: AlbumId): Promise<Media[]> {
        if (this.shouldFailFetch) {
            throw new Error("Failed to fetch medias");
        }
        return this.medias;
    }
}

describe("thunk:updateAlbumDates", () => {
    it("should convert display dates to API format with default times and dispatch actions", async () => {
        const rawMedias = someMedias.flatMap(m => m.medias);
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, rawMedias);
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            {
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate,
                startAtDayStart: true,
                endAtDayEnd: true,
            }
        );

        expect(fakePort.updatedAlbums).toEqual([{
            albumId,
            startDate: new Date("2023-07-10T00:00:00.000Z"),
            endDate: new Date("2023-07-21T00:00:00.000Z"),
        }]);

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdated({
                albums: twoAlbums,
                medias: someMedias,
            })
        ]);
    });

    it("should convert specific times to API format with exclusive end", async () => {
        const rawMedias = someMedias.flatMap(m => m.medias);
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, rawMedias);
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10T10:30:00");
        const displayEndDate = new Date("2023-07-20T15:00:00");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            {
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate,
                startAtDayStart: false,
                endAtDayEnd: false,
            }
        );

        expect(fakePort.updatedAlbums).toEqual([{
            albumId,
            startDate: new Date("2023-07-10T10:30:00.000Z"),
            endDate: new Date("2023-07-20T15:01:00.000Z"),
        }]);

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdated({
                albums: twoAlbums,
                medias: someMedias,
            })
        ]);
    });

    it("should prioritize start/end of day flags over specific times", async () => {
        const rawMedias = someMedias.flatMap(m => m.medias);
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, rawMedias);
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10T12:34:56.789Z");
        const displayEndDate = new Date("2023-07-20T23:45:01.234Z");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            {
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate,
                startAtDayStart: true,
                endAtDayEnd: true,
            }
        );

        expect(fakePort.updatedAlbums).toEqual([{
            albumId,
            startDate: new Date("2023-07-10T00:00:00.000Z"),
            endDate: new Date("2023-07-21T00:00:00.000Z"),
        }]);

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdated({
                albums: twoAlbums,
                medias: someMedias,
            })
        ]);
    });

    it("should dispatch failure action when update fails", async () => {
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, []);
        fakePort.shouldFailUpdate = true;
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            {
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate,
                startAtDayStart: true,
                endAtDayEnd: true,
            }
        );

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdateFailed({error: "Failed to update album dates"})
        ]);
    });

    it("should dispatch failure action when fetch fails", async () => {
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, []);
        fakePort.shouldFailFetch = true;
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            {
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate,
                startAtDayStart: true,
                endAtDayEnd: true,
            }
        );

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdateFailed({error: "Failed to fetch albums"})
        ]);
    });
});
