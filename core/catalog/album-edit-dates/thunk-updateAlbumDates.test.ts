import {updateAlbumDatesThunk, UpdateAlbumDatesPort} from "./thunk-updateAlbumDates";
import {EditDatesDialogState} from "../language/CatalogViewerState";
import {loadingStarted} from "./action-loadingStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {twoAlbums, twoMedias} from "../tests/test-helper-state";

class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
    updatedAlbums: { albumId: AlbumId; startDate: Date; endDate: Date }[] = [];

    async updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void> {
        this.updatedAlbums.push({albumId, startDate, endDate});
    }

    async fetchAlbums(): Promise<Album[]> {
        return twoAlbums;
    }

    async fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return twoMedias;
    }
}

describe("thunk:updateAlbumDates", () => {
    const editDatesDialogState: EditDatesDialogState = {
        albumId: {owner: "myself", folderName: "album1"},
        albumName: "Test Album",
        startDate: new Date("2023-07-10T00:00:00"),
        endDate: new Date("2023-07-21T00:00:00"),
        isLoading: false
    };

    it("should update album dates and refresh data", async () => {
        const fakePort = new UpdateAlbumDatesPortFake();
        const dispatched: any[] = [];

        await updateAlbumDatesThunk(
            dispatched.push.bind(dispatched),
            fakePort,
            editDatesDialogState
        );

        expect(fakePort.updatedAlbums).toEqual([{
            albumId: {owner: "myself", folderName: "album1"},
            startDate: new Date("2023-07-10T00:00:00"),
            endDate: new Date("2023-07-21T00:00:00")
        }]);

        expect(dispatched).toEqual([
            loadingStarted(),
            albumDatesUpdated({
                albums: twoAlbums,
                medias: twoMedias,
                updatedAlbumId: {owner: "myself", folderName: "album1"}
            })
        ]);
    });
});
