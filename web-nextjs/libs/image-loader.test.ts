import {describe, expect, it} from 'vitest';
import {catalogImageLoader} from './image-loader';

describe('catalogImageLoader', () => {
    it('appends the requested width as a query parameter to the content path', () => {
        const result = catalogImageLoader({src: '/api/v1/owners/myself/medias/media-1/image.jpg', width: 1440, quality: 75});

        expect(result).toBe('/api/v1/owners/myself/medias/media-1/image.jpg?w=1440');
    });

    it('appends the requested width to a plain content path', () => {
        const result = catalogImageLoader({src: '/content/foo.jpg', width: 2400, quality: 75});

        expect(result).toBe('/content/foo.jpg?w=2400');
    });

    it('caps the requested width at the maximum cacheable resolution', () => {
        const result = catalogImageLoader({src: '/content/foo.jpg', width: 3840, quality: 75});

        expect(result).toBe('/content/foo.jpg?w=2400');
    });
});
