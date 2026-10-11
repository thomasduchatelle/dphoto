import type {ImageLoaderProps} from 'next/image';
import {mediaUrl} from './requests/media-url';

export function catalogImageLoader({src, width}: ImageLoaderProps): string {
    return mediaUrl(src, width) ?? src;
}
