import {updateAlbumDatesThunk, UpdateAlbumDatesPort} from "./thunk-updateAlbumDates";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {twoAlbums, someMedias} from "../tests/test-helper-state";
import {Album, Media, AlbumId, MediaWithinADay} from "../language";
import {MediaPerDayLoader} from "../navigation/MediaPerDayLoader"; // Import MediaPerDayLoader

class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
    public updatedAlbums: {albumId: AlbumId, startDate: Date, endDate: Date}[] = [];
    
    constructor(
        private albums: Album[] = [],
        private medias: Media[] = [] // Changed to Media[]
    ) {}

    async updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void> {
        this.updatedAlbums.push({albumId, startDate, endDate});
    }

    async fetchAlbums(): Promise<Album[]> {
        return this.albums;
    }

    async fetchMedias(albumId: AlbumId): Promise<Media[]> { // Changed to Media[]
        return this.medias;
    }
}

// Mock MediaPerDayLoader
class MediaPerDayLoaderFake extends MediaPerDayLoader {
    constructor(private mockMedias: MediaWithinADay[]) {
        super(null as any); // Pass null or a mock for FetchAlbumMediasPort as it's not used in this mock
    }

    public async loadMedias(albumId: AlbumId): Promise<{ medias: MediaWithinADay[] }> {
        return { medias: this.mockMedias };
    }
}

describe("thunk:updateAlbumDates", () => {
    it("should convert display dates to API format and dispatch actions", async () => {
        const fakePort = new UpdateAlbumDatesPortFake(twoAlbums, someMedias[0].medias); // Pass raw medias
        const mediaPerDayLoaderFake = new MediaPerDayLoaderFake(someMedias); // Pass grouped medias
        const dispatched: any[] = [];
        const albumId = twoAlbums[0].albumId;
        const displayStartDate = new Date("2023-07-10");
        const displayEndDate = new Date("2023-07-20");

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            mediaPerDayLoaderFake, // Pass the fake MediaPerDayLoader
            albumId,
            displayStartDate,
            displayEndDate
        );

        expect(fakePort.updatedAlbums).toEqual([{
            albumId,
            startDate: new Date(Date.UTC(2023, 6, 10, 0, 0, 0, 0)), // Expected UTC date for 2023-07-10 00:00:00 local
            endDate: new Date(Date.UTC(2023, 6, 21, 0, 0, 0, 0)),   // Expected UTC date for 2023-07-21 00:00:00 local
        }]);

        expect(dispatched).toEqual([
            albumDatesUpdateStarted(),
            albumDatesUpdated({
                albums: twoAlbums,
                medias: someMedias,
            })
        ]);
    });
});
