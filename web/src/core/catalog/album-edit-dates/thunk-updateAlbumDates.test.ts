import {UpdateAlbumDatesPort, updateAlbumDatesThunk} from "./thunk-updateAlbumDates";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {someMedias, twoAlbums} from "../tests/test-helper-state";
import {Album, AlbumId, Media} from "../language";
import {groupByDay} from "../navigation/group-by-day";

class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
    public updatedAlbums: { albumId: AlbumId, startDate: Date, endDate: Date }[] = [];

    constructor(
        private albums: Album[] = [],
        private medias: Media[] = []
    ) {
    }

    async updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void> {
        this.updatedAlbums.push({albumId, startDate, endDate});
    }

    async fetchAlbums(): Promise<Album[]> {
        return this.albums;
    }

    async fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return this.medias;
    }
}

describe("thunk:updateAlbumDates", () => {
    it("should convert display dates to API format and dispatch actions", async () => {
        const rawMedias = someMedias.flatMap(m => m.medias);
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, rawMedias);
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            { // Pass arguments as a single object
                albumId,
                startDate: displayStartDate,
                endDate: displayEndDate
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
                medias: groupByDay(rawMedias),
            })
        ]);
    });
});
