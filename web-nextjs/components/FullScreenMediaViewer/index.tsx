'use client';

import Lightbox, {Slide} from 'yet-another-react-lightbox';
import Counter from 'yet-another-react-lightbox/plugins/counter';
import Download from 'yet-another-react-lightbox/plugins/download';
import Slideshow from 'yet-another-react-lightbox/plugins/slideshow';
import Video from 'yet-another-react-lightbox/plugins/video';
import Zoom from 'yet-another-react-lightbox/plugins/zoom';
import 'yet-another-react-lightbox/styles.css';
import 'yet-another-react-lightbox/plugins/counter.css';
import {Media, MediaType} from '@/domains/catalog/language';
import {mediaUrl, prefixRelativeUrl} from '@/libs/requests/media-url';
import {runtimeApiPrefix} from '@/libs/requests/basepath';

export interface FullScreenMediaViewerProps {
    medias: Media[];
    displayedIndex: number;
    slideshow?: boolean;
    onClose: () => void;
    onNavigate: (mediaId: string) => void;
}

// Widths requested through the custom image URL (`contentPath?w=<width>`). YARL
// picks the smallest source the viewport needs and upgrades as larger ones
// arrive — so the already-cached thumbnail shows instantly while the full-res
// image loads. The full-res sources also give the Zoom plugin the real image
// to magnify (the plugin stays disabled when no native image slide is rendered).
const FULL_RES_WIDTHS = [1440, 2400] as const;
// Approximate rendered width of the grid thumbnail; only used so the browser
// ranks the thumbnail below the full-res sources in the srcSet.
const THUMBNAIL_WIDTH = 360;

function toSlide(media: Media): Slide {
    const originalUrl = prefixRelativeUrl(media.contentPath, runtimeApiPrefix) ?? media.contentPath;
    const filename = media.contentPath.split('/').pop() || media.id;
    const download = {url: originalUrl, filename};

    if (media.type === MediaType.VIDEO) {
        return {
            type: 'video',
            poster: media.thumbnailUrl,
            sources: [{src: originalUrl, type: 'video/mp4'}],
            download,
        };
    }

    const srcSet = [
        {src: media.thumbnailUrl, width: THUMBNAIL_WIDTH, height: THUMBNAIL_WIDTH},
        ...FULL_RES_WIDTHS.map(width => ({src: mediaUrl(media.contentPath, width), width, height: width})),
    ];

    return {
        type: 'image',
        src: mediaUrl(media.contentPath, FULL_RES_WIDTHS[FULL_RES_WIDTHS.length - 1]),
        srcSet,
        download,
    };
}

export function FullScreenMediaViewer({medias, displayedIndex, slideshow, onClose, onNavigate}: FullScreenMediaViewerProps) {
    return (
        <Lightbox
            open
            close={onClose}
            index={displayedIndex}
            slides={medias.map(toSlide)}
            plugins={[Counter, Download, Slideshow, Video, Zoom]}
            counter={{separator: 'of'}}
            carousel={{finite: true, padding: 0}}
            controller={{closeOnBackdropClick: true}}
            slideshow={{autoplay: slideshow}}
            on={{
                view: ({index}) => {
                    const media = medias[index];
                    if (media && index !== displayedIndex) {
                        onNavigate(media.id);
                    }
                },
            }}
        />
    );
}
