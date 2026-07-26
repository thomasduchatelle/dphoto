import {albumKey, AlbumTemperatureInput, computeAlbumTemperatures} from "./utils-computeAlbumTemperatures";
import {AlbumId} from "./catalog-state";

const albumId = (folderName: string): AlbumId => ({owner: "tomdush@gmail.com", folderName});

const albumOverNDays = (folderName: string, totalCount: number, days: number): AlbumTemperatureInput => ({
    albumId: albumId(folderName),
    start: new Date("2024-01-01T00:00:00Z"),
    end: new Date(new Date("2024-01-01T00:00:00Z").getTime() + days * 24 * 3600 * 1000),
    totalCount,
});

describe("utils:computeAlbumTemperatures", () => {
    it("should return an empty map when there are no albums", () => {
        const result = computeAlbumTemperatures([]);

        expect(result.size).toBe(0);
    });

    it("should assign a relativeTemperature of 0.5 to a single album", () => {
        const result = computeAlbumTemperatures([albumOverNDays("only", 100, 10)]);

        expect(result.get(albumKey(albumId("only")))).toEqual({temperature: 10, relativeTemperature: 0.5});
    });

    it("should rank albums by temperature and spread relativeTemperature evenly across [0, 1]", () => {
        const result = computeAlbumTemperatures([
            albumOverNDays("coldest", 10, 10), // temperature 1
            albumOverNDays("cold", 20, 10), // temperature 2
            albumOverNDays("hot", 30, 10), // temperature 3
            albumOverNDays("hottest", 40, 10), // temperature 4
        ]);

        expect(result.get(albumKey(albumId("coldest")))).toEqual({temperature: 1, relativeTemperature: 0.125});
        expect(result.get(albumKey(albumId("cold")))).toEqual({temperature: 2, relativeTemperature: 0.375});
        expect(result.get(albumKey(albumId("hot")))).toEqual({temperature: 3, relativeTemperature: 0.625});
        expect(result.get(albumKey(albumId("hottest")))).toEqual({temperature: 4, relativeTemperature: 0.875});
    });

    it("should assign the same relativeTemperature to albums with the same temperature", () => {
        const result = computeAlbumTemperatures([
            albumOverNDays("first-tied", 10, 10), // temperature 1
            albumOverNDays("second-tied", 20, 20), // temperature 1
            albumOverNDays("hottest", 40, 10), // temperature 4
        ]);

        expect(result.get(albumKey(albumId("first-tied")))?.relativeTemperature)
            .toBe(result.get(albumKey(albumId("second-tied")))?.relativeTemperature);
        expect(result.get(albumKey(albumId("hottest")))?.relativeTemperature).toBeCloseTo(5 / 6);
    });

    it("should not divide by zero when start and end are the same instant", () => {
        const sameInstant = new Date("2024-01-01T00:00:00Z");

        const result = computeAlbumTemperatures([{
            albumId: albumId("instant"),
            start: sameInstant,
            end: sameInstant,
            totalCount: 5,
        }]);

        expect(result.get(albumKey(albumId("instant")))?.temperature).toBe(5);
    });

    it("should distribute a skewed dataset into ~25% quartile buckets", () => {
        const albums: AlbumTemperatureInput[] = Array.from({length: 200}, (_, i) => {
            const totalCount = i < 10 ? (i + 1) * 500 : (i % 50) + 1;
            return albumOverNDays(`album-${i}`, totalCount, 10);
        });

        const result = computeAlbumTemperatures(albums);
        const buckets = [0, 0, 0, 0];
        for (const {relativeTemperature} of result.values()) {
            if (relativeTemperature >= 0.75) buckets[3]++;
            else if (relativeTemperature >= 0.5) buckets[2]++;
            else if (relativeTemperature >= 0.25) buckets[1]++;
            else buckets[0]++;
        }

        buckets.forEach(bucketSize => {
            expect(bucketSize).toBeGreaterThanOrEqual(45);
            expect(bucketSize).toBeLessThanOrEqual(55);
        });
    });
});
