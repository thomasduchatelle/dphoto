'use client';

import {useState} from 'react';
import Image from 'next/image';
import {Box} from '@mui/material';
import BrokenImageOutlinedIcon from '@mui/icons-material/BrokenImageOutlined';
import Lightbox, {RenderSlideProps, Slide, SlideImage} from 'yet-another-react-lightbox';
import Counter from 'yet-another-react-lightbox/plugins/counter';
import Download from 'yet-another-react-lightbox/plugins/download';
import Video from 'yet-another-react-lightbox/plugins/video';
import Zoom from 'yet-another-react-lightbox/plugins/zoom';
import 'yet-another-react-lightbox/styles.css';
import 'yet-another-react-lightbox/plugins/counter.css';
import {Media, MediaType} from '@/domains/catalog/language';
import {catalogImageLoader} from '@/libs/image-loader';
import {prefixRelativeUrl} from '@/libs/requests/media-url';
import {runtimeApiPrefix} from '@/libs/requests/basepath';

export interface FullScreenMediaViewerProps {
    medias: Media[];
    displayedIndex: number;
    onClose: () => void;
    onNavigate: (mediaId: string) => void;
}

function MediaSlide({slide, rect}: RenderSlideProps<SlideImage>) {
    const [failed, setFailed] = useState(false);

    if (failed) {
        return (
            <Box
                role="img"
                aria-label="This media could not be loaded"
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: 1,
                    color: 'rgba(255,255,255,0.5)',
                }}
            >
                <BrokenImageOutlinedIcon sx={{fontSize: 64}}/>
                <Box component="span" sx={{fontSize: '0.9rem'}}>This media could not be loaded</Box>
            </Box>
        );
    }

    const sizes = `(max-width: 900px) 1440px, 2400px`;
    const width = rect.width > 0 ? Math.round(rect.width) : 1440;
    const height = rect.height > 0 ? Math.round(rect.height) : 1440;

    return (
        <Image
            loader={catalogImageLoader}
            src={slide.src}
            alt=""
            width={width}
            height={height}
            sizes={sizes}
            style={{maxWidth: '100%', maxHeight: '100%', width: 'auto', height: 'auto', objectFit: 'contain'}}
            onError={() => setFailed(true)}
            priority
        />
    );
}

function renderSlide(props: RenderSlideProps) {
    if (props.slide.type !== 'image') {
        return undefined;
    }
    return <MediaSlide key={props.slide.src} {...(props as RenderSlideProps<SlideImage>)}/>;
}

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

    return {
        type: 'image',
        src: media.contentPath,
        download,
    };
}

export function FullScreenMediaViewer({medias, displayedIndex, onClose, onNavigate}: FullScreenMediaViewerProps) {
    return (
        <Lightbox
            open
            close={onClose}
            index={displayedIndex}
            slides={medias.map(toSlide)}
            plugins={[Counter, Download, Video, Zoom]}
            counter={{separator: 'of'}}
            carousel={{finite: true, padding: 0}}
            controller={{closeOnBackdropClick: true}}
            on={{
                view: ({index}) => {
                    const media = medias[index];
                    if (media && index !== displayedIndex) {
                        onNavigate(media.id);
                    }
                },
            }}
            render={{slide: renderSlide}}
        />
    );
}
