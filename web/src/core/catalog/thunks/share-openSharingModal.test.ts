import {openSharingModalThunk} from "./share-openSharingModal";
import {Album, catalogActions} from "../domain";

describe("openSharingModalThunk", () => {
    it("should dispatch openSharingModalAction with albumId when called", () => {
        const album: Album = {
            albumId: {owner: "owner", folderName: "album"},
            name: "Test Album",
            start: new Date(),
            end: new Date(),
            totalCount: 0,
            temperature: 0,
            relativeTemperature: 0,
            sharedWith: [],
        };
        const dispatched: any[] = [];
        openSharingModalThunk(dispatched.push.bind(dispatched), album);
        expect(dispatched).toEqual([
            catalogActions.openSharingModalAction({albumId: album.albumId})
        ]);
    });
});
