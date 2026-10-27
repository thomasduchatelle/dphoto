import {describe, expect, it} from 'vitest';
import {albumRailWindow} from './album-rail-window';
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

const names = (albums: Album[]): string[] => albums.map((album) => album.name);

describe('albumRailWindow', () => {
    const tenAlbums = Array.from({length: 10}, (_, index) => albumAt(index));

    it('returns the whole list when it has fewer than 7 albums', () => {
        const albums = tenAlbums.slice(0, 5);

        expect(albumRailWindow(albums, 2)).toEqual(albums);
    });

    it('returns the whole list when it has exactly 7 albums', () => {
        const albums = tenAlbums.slice(0, 7);

        expect(albumRailWindow(albums, 4)).toEqual(albums);
    });

    it('centres the window with 3 before and 3 after the displayed album', () => {
        expect(names(albumRailWindow(tenAlbums, 5))).toEqual([
            'Album 2', 'Album 3', 'Album 4', 'Album 5', 'Album 6', 'Album 7', 'Album 8',
        ]);
    });

    it('shows the displayed album and the 6 following when it is first', () => {
        expect(names(albumRailWindow(tenAlbums, 0))).toEqual([
            'Album 0', 'Album 1', 'Album 2', 'Album 3', 'Album 4', 'Album 5', 'Album 6',
        ]);
    });

    it('shows the 6 preceding albums and the displayed one when it is last', () => {
        expect(names(albumRailWindow(tenAlbums, 9))).toEqual([
            'Album 3', 'Album 4', 'Album 5', 'Album 6', 'Album 7', 'Album 8', 'Album 9',
        ]);
    });

    it('biases the window to keep 7 albums when the displayed album is near the end', () => {
        expect(names(albumRailWindow(tenAlbums, 8))).toEqual([
            'Album 3', 'Album 4', 'Album 5', 'Album 6', 'Album 7', 'Album 8', 'Album 9',
        ]);
    });

    it('returns the whole list when the displayed album is not found', () => {
        expect(albumRailWindow(tenAlbums, -1)).toEqual(tenAlbums);
    });
});
