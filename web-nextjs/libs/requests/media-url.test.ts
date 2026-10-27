import {describe, expect, it} from 'vitest';
import {mediaUrl} from './media-url';

describe('mediaUrl', () => {
    it('appends the requested width to the media url', () => {
        expect(mediaUrl('/api/v1/owners/me/medias/1/media.jpg', 360)).toBe('/api/v1/owners/me/medias/1/media.jpg?w=360');
    });
});
