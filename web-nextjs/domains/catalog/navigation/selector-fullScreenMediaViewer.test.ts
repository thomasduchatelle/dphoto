import {describe, expect, it} from 'vitest';
import {selectFullScreenMediaViewer} from './selector-fullScreenMediaViewer';
import {loadedStateWithTwoAlbums, someMediasByDays, twoAlbums} from '../tests/test-helper-state';
import {CatalogViewerState, Media, MediaType} from '../language';

const mediaA: Media = {
    id: 'media-a',
    type: MediaType.IMAGE,
    time: new Date('2025-01-05T12:00:00Z'),
    uiRelativePath: 'media-a/image.jpg',
    contentPath: '/media-a.jpg',
    thumbnailUrl: '/media-a.jpg?w=360',
    source: '',
};

const mediaB: Media = {
    id: 'media-b',
    type: MediaType.IMAGE,
    time: new Date('2025-01-04T12:00:00Z'),
    uiRelativePath: 'media-b/image.jpg',
    contentPath: '/media-b.jpg',
    thumbnailUrl: '/media-b.jpg?w=360',
    source: '',
};

const mediaC: Media = {
    id: 'media-c',
    type: MediaType.IMAGE,
    time: new Date('2025-01-03T12:00:00Z'),
    uiRelativePath: 'media-c/image.jpg',
    contentPath: '/media-c.jpg',
    thumbnailUrl: '/media-c.jpg?w=360',
    source: '',
};

const stateWithMultipleMedias: CatalogViewerState = {
    ...loadedStateWithTwoAlbums,
    medias: [
        {day: new Date(2025, 0, 5), medias: [mediaA, mediaB]},
        {day: new Date(2025, 0, 3), medias: [mediaC]},
    ],
};

describe('selectFullScreenMediaViewer', () => {
    it('returns a flat media list, the displayed index, and the album id when the media is in the album', () => {
        const result = selectFullScreenMediaViewer(stateWithMultipleMedias, 'media-b');

        expect(result).toEqual({
            albumId: twoAlbums[0].albumId,
            medias: [mediaA, mediaB, mediaC],
            displayedIndex: 1,
            mediaNotFound: false,
        });
    });

    it('flags the media as not found when the mediaId is not part of the album', () => {
        const result = selectFullScreenMediaViewer(loadedStateWithTwoAlbums, 'unknown-media');

        expect(result).toEqual({
            albumId: twoAlbums[0].albumId,
            medias: someMediasByDays.flatMap(group => group.medias),
            displayedIndex: -1,
            mediaNotFound: true,
        });
    });
});
