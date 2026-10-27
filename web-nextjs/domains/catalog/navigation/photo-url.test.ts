import {describe, expect, it} from 'vitest';
import {photoUrl} from './photo-url';

describe('photoUrl', () => {
    it('builds the full photo URL from an album id and a media id', () => {
        expect(photoUrl({owner: 'myself', folderName: 'jan-25'}, 'media-1')).toBe('/albums/myself/jan-25/photos/media-1');
    });

    it('url-encodes special characters in the owner, folder name, and media id', () => {
        expect(photoUrl({owner: 'me/you', folderName: 'jan 25'}, 'id with spaces')).toBe('/albums/me%2Fyou/jan%2025/photos/id%20with%20spaces');
    });

    it('appends the play flag to start the slideshow when autoplay is requested', () => {
        expect(photoUrl({owner: 'myself', folderName: 'jan-25'}, 'media-1', {autoplay: true})).toBe('/albums/myself/jan-25/photos/media-1?play');
    });
});
