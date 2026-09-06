import {describe, it, expect, beforeEach} from 'vitest';
import {Album, AlbumId, Media, MediaType} from '../language';
import {albumsAndMediasLoaded} from './action-albumsAndMediasLoaded';
import {mediaLoadFailed} from './action-mediaLoadFailed';
import {noAlbumAvailable} from './action-noAlbumAvailable';
import {LoadAlbumPage, LoadAlbumPagePort} from './thunk-loadAlbumPage';
import {twoAlbums} from '../tests/test-helper-state';
import {CatalogViewerAction} from '../actions';

const albumId = twoAlbums[0].albumId;

const someMedias: Media[] = [
    {
        id: 'media-1',
        type: MediaType.IMAGE,
        time: new Date('2025-01-05T12:42:00Z'),
        uiRelativePath: 'media-1/image.jpg',
        contentPath: '/media-1.jpg',
        source: '',
    },
];

class FakeLoadAlbumPageAdapter implements LoadAlbumPagePort {
    constructor(
        private readonly albums: Album[],
        private readonly mediasByAlbum: Map<string, Media[]> = new Map(),
    ) {}

    fetchAlbums(): Promise<Album[]> {
        return Promise.resolve(this.albums);
    }

    fetchMedias(id: AlbumId): Promise<Media[]> {
        return Promise.resolve(this.mediasByAlbum.get(`${id.owner}/${id.folderName}`) ?? []);
    }
}

class FakeAdapterWithAlbumFailure implements LoadAlbumPagePort {
    constructor(private readonly error: Error) {}

    fetchAlbums(): Promise<Album[]> {
        return Promise.reject(this.error);
    }

    fetchMedias(): Promise<Media[]> {
        return Promise.resolve([]);
    }
}

class FakeAdapterWithMediaFailure implements LoadAlbumPagePort {
    constructor(private readonly albums: Album[], private readonly error: Error) {}

    fetchAlbums(): Promise<Album[]> {
        return Promise.resolve(this.albums);
    }

    fetchMedias(): Promise<Media[]> {
        return Promise.reject(this.error);
    }
}

describe('LoadAlbumPage', () => {
    let dispatched: CatalogViewerAction[];

    const newLoader = (port: LoadAlbumPagePort) =>
        new LoadAlbumPage((a) => dispatched.push(a as CatalogViewerAction), port);

    beforeEach(() => {
        dispatched = [];
    });

    it('should dispatch albumsAndMediasLoaded when the album exists and medias load successfully', async () => {
        const port = new FakeLoadAlbumPageAdapter(
            twoAlbums,
            new Map([[`${albumId.owner}/${albumId.folderName}`, someMedias]]),
        );

        await newLoader(port).loadAlbumPage(albumId);

        expect(dispatched).toEqual([
            albumsAndMediasLoaded({albums: twoAlbums, medias: someMedias, mediasFromAlbumId: albumId}),
        ]);
    });

    it('should dispatch noAlbumAvailable when the albumId is not in the fetched album list', async () => {
        const port = new FakeLoadAlbumPageAdapter(twoAlbums);
        const unknownId: AlbumId = {owner: 'nobody', folderName: 'ghost'};

        await newLoader(port).loadAlbumPage(unknownId);

        expect(dispatched).toEqual([noAlbumAvailable(undefined)]);
    });

    it('should dispatch mediaLoadFailed when medias cannot be fetched', async () => {
        const error = new Error('media fetch failed');
        const port = new FakeAdapterWithMediaFailure(twoAlbums, error);

        await newLoader(port).loadAlbumPage(albumId);

        expect(dispatched).toEqual([
            mediaLoadFailed({
                albums: twoAlbums,
                displayedAlbumId: albumId,
                error: new Error(`failed to load medias of ${JSON.stringify(albumId)}`, {cause: error}),
            }),
        ]);
    });

    it('should reject when albums cannot be fetched', async () => {
        const error = new Error('albums fetch failed');
        const port = new FakeAdapterWithAlbumFailure(error);

        await expect(newLoader(port).loadAlbumPage(albumId)).rejects.toThrow(error);
        expect(dispatched).toHaveLength(0);
    });
});
