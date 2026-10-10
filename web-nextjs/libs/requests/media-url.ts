import {basePath, runtimeApiPrefix} from "./basepath";

export function prefixRelativeUrl(url: string | undefined, prefix: string): string | undefined {
    if (!url || !prefix) return url;
    if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith(prefix)) return url;
    return `${prefix}${url}`;
}

export function withBasePath(url: string): string {
    return prefixRelativeUrl(url, basePath) ?? url;
}

export function mediaUrl(contentPath: string, width: number, prefix: string = runtimeApiPrefix): string {
    return `${prefixRelativeUrl(contentPath, prefix)}?w=${width}`;
}
