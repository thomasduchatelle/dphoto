import {updateAlbumDatesThunk, UpdateAlbumDatesPort} from "./thunk-updateAlbumDates";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {twoAlbums, someMedias} from "../tests/test-helper-state";
import {Album, Media, AlbumId} from "../language";

class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
    public updatedAlbums: {albumId: AlbumId, startDate: Date, endDate: Date}[] = [];
    
    constructor(
        private albums: Album[] = [],
        private medias: Media[] = []
    ) {}

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
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, someMedias[0].medias);
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            albumId,
            displayStartDate,
            displayEndDate
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
                medias: someMedias[0].medias,
                updatedAlbumId: albumId,
            })
        ]);
    });
});
