'use client';

import Lightbox, {IconButton, Label, Slide, useLightboxState} from 'yet-another-react-lightbox';
import Counter from 'yet-another-react-lightbox/plugins/counter';
import Download from 'yet-another-react-lightbox/plugins/download';
import Slideshow from 'yet-another-react-lightbox/plugins/slideshow';
import Video from 'yet-another-react-lightbox/plugins/video';
import Zoom from 'yet-another-react-lightbox/plugins/zoom';
import 'yet-another-react-lightbox/styles.css';
import 'yet-another-react-lightbox/plugins/counter.css';
import StarIcon from '@mui/icons-material/Star';
import StarBorderIcon from '@mui/icons-material/StarBorder';
import {FAVORITE_TAG, Media, MediaId, MediaType} from '@/domains/catalog/language';
import {mediaUrl} from '@/libs/requests/media-url';

export interface FullScreenMediaViewerProps {
    medias: Media[];
    displayedIndex: number;
    slideshow?: boolean;
    isOwner: boolean;
    onClose: () => void;
    onNavigate: (mediaId: string) => void;
    onAddTag: (mediaId: MediaId, tag: string) => void;
    onRemoveTag: (mediaId: MediaId, tag: string) => void;
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
    const originalUrl = media.contentPath;
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

function FavoriteButton({medias, isOwner, onAddTag, onRemoveTag}: Pick<FullScreenMediaViewerProps, 'medias' | 'isOwner' | 'onAddTag' | 'onRemoveTag'>) {
    const {currentIndex} = useLightboxState();
    const media = medias[currentIndex];
    if (!media) {
        return null;
    }

    const isFavorite = media.tags.includes(FAVORITE_TAG);

    return (
        <IconButton
            key="favorite"
            label={(isFavorite ? 'Remove from favourites' : 'Mark as favourite') as Label}
            icon={isFavorite ? StarIcon : StarBorderIcon}
            disabled={!isOwner}
            onClick={() => isFavorite ? onRemoveTag(media.id, FAVORITE_TAG) : onAddTag(media.id, FAVORITE_TAG)}
        />
    );
}

export function FullScreenMediaViewer({medias, displayedIndex, slideshow, isOwner, onClose, onNavigate, onAddTag, onRemoveTag}: FullScreenMediaViewerProps) {
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
            toolbar={{
                buttons: [
                    <FavoriteButton key="favorite" medias={medias} isOwner={isOwner} onAddTag={onAddTag} onRemoveTag={onRemoveTag}/>,
                    'download',
                    'slideshow',
                    'close',
                ],
            }}
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
