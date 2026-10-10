import {describe, expect, it} from 'vitest';
import {mediaUrl, prefixRelativeUrl, withBasePath} from './media-url';
import {basePath} from './basepath';

describe('prefixRelativeUrl', () => {
    it('prepends the prefix to a relative url', () => {
        expect(prefixRelativeUrl('/api/v1/medias/1', '/nextjs')).toBe('/nextjs/api/v1/medias/1');
    });

    it('does not prepend the prefix when the url already starts with it', () => {
        expect(prefixRelativeUrl('/nextjs/thumbnails/1.jpg', '/nextjs')).toBe('/nextjs/thumbnails/1.jpg');
    });

    it('leaves absolute urls and undefined untouched', () => {
        expect(prefixRelativeUrl('https://cdn.example.com/1.jpg', '/nextjs')).toBe('https://cdn.example.com/1.jpg');
        expect(prefixRelativeUrl(undefined, '/nextjs')).toBeUndefined();
    });

    it('is a no-op without a prefix', () => {
        expect(prefixRelativeUrl('/api/v1/medias/1', '')).toBe('/api/v1/medias/1');
    });
});

describe('withBasePath', () => {
    it('prefixes a static asset with the application basePath', () => {
        expect(withBasePath('/video-placeholder.png')).toBe(`${basePath}/video-placeholder.png`);
    });

    it('is idempotent', () => {
        expect(withBasePath(`${basePath}/video-placeholder.png`)).toBe(`${basePath}/video-placeholder.png`);
    });
});

describe('mediaUrl', () => {
    it('appends the requested width to the media url', () => {
        expect(mediaUrl('/api/v1/owners/me/medias/1/media.jpg', 360, '/nextjs')).toBe('/nextjs/api/v1/owners/me/medias/1/media.jpg?w=360');
    });

    it('does not prefix the media url when the runtime does not require it', () => {
        expect(mediaUrl('/api/v1/owners/me/medias/1/media.jpg', 257, '')).toBe('/api/v1/owners/me/medias/1/media.jpg?w=257');
    });
});
