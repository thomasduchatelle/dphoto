import {describe, expect, it} from 'vitest';
import {catalogViewerPageSelector} from './selector-catalog-viewer-page';
import {loadedStateWithTwoAlbums, march2025, twoAlbums} from '../tests/test-helper-state';

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
});
