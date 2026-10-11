import type {ImageLoaderProps} from 'next/image';
import {mediaUrl} from './requests/media-url';

const MAX_CACHEABLE_WIDTH = 2400;

export function catalogImageLoader({src, width}: ImageLoaderProps): string {
    return mediaUrl(src, Math.min(width, MAX_CACHEABLE_WIDTH)) ?? src;
}
