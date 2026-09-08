import {describe, expect, it} from 'vitest';
import {albumUrl} from './album-url';

describe('albumUrl', () => {
    it('builds the album page path from the album id', () => {
        expect(albumUrl({owner: 'sandfall', folderName: 'clair-obscur'})).toBe('/albums/sandfall/clair-obscur');
    });

    it('URL-encodes the owner and the folder name', () => {
        expect(albumUrl({owner: 'tony@stark.com', folderName: 'summer 2026'})).toBe('/albums/tony%40stark.com/summer%202026');
    });
});
