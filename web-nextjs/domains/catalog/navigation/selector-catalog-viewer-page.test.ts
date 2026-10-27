import {describe, expect, it} from 'vitest';
import {catalogViewerPageSelector} from './selector-catalog-viewer-page';
import {loadedStateWithTwoAlbums, march2025, twoAlbums} from '../tests/test-helper-state';
import {Album} from '../language';

const albumAt = (index: number): Album => ({
    albumId: {owner: 'myself', folderName: `album-${index}`},
    name: `Album ${index}`,
    start: new Date(2025, 0, 1),
    end: new Date(2025, 1, 1),
    totalCount: 0,
    temperature: 0,
    relativeTemperature: 0,
    sharedWith: [],
    covers: [],
});

describe('catalogViewerPageSelector', () => {
    it('returns nextAlbum (newer, index-1) and previousAlbum (older, index+1) relative to the displayed album', () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0], twoAlbums[1], march2025],
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
        };

        const result = catalogViewerPageSelector(state);

        expect(result.nextAlbum).toEqual(twoAlbums[0]);
        expect(result.previousAlbum).toEqual(march2025);
    });

    it('returns undefined for both when the album is alone', () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
        };

        const result = catalogViewerPageSelector(state);

        expect(result.previousAlbum).toBeUndefined();
        expect(result.nextAlbum).toBeUndefined();
    });

    it('returns the photo URL of the first media as firstMediaHref', () => {
        const result = catalogViewerPageSelector(loadedStateWithTwoAlbums);

        expect(result.firstMediaHref).toEqual('/albums/myself/jan-25/photos/media-1');
    });

    it('returns undefined firstMediaHref when the album has no media', () => {
        const result = catalogViewerPageSelector({...loadedStateWithTwoAlbums, medias: []});

        expect(result.firstMediaHref).toBeUndefined();
    });

    it('limits railAlbums to a 7-album window centred on the displayed album within the filtered list', () => {
        const albums = Array.from({length: 10}, (_, index) => albumAt(index));
        const state = {
            ...loadedStateWithTwoAlbums,
            allAlbums: albums,
            albums,
            mediasLoadedFromAlbumId: albums[5].albumId,
        };

        const result = catalogViewerPageSelector(state);

        expect(result.railAlbums.map(album => album.name)).toEqual([
            'Album 2', 'Album 3', 'Album 4', 'Album 5', 'Album 6', 'Album 7', 'Album 8',
        ]);
    });
});
