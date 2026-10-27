'use client';

import {useState} from 'react';
import Image from 'next/image';
import {Box} from '@mui/material';
import BrokenImageOutlinedIcon from '@mui/icons-material/BrokenImageOutlined';
import Lightbox, {RenderSlideProps, SlideImage} from 'yet-another-react-lightbox';
import Slideshow from 'yet-another-react-lightbox/plugins/slideshow';
import 'yet-another-react-lightbox/styles.css';
import {Media} from '@/domains/catalog/language';
import {catalogImageLoader} from '@/libs/image-loader';

export interface FullScreenMediaViewerProps {
    medias: Media[];
    displayedIndex: number;
    slideshow?: boolean;
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

function renderSlide(props: RenderSlideProps<SlideImage>) {
    return <MediaSlide key={props.slide.src} {...props}/>;
}

export function FullScreenMediaViewer({medias, displayedIndex, slideshow, onClose, onNavigate}: FullScreenMediaViewerProps) {
    const slides: SlideImage[] = medias.map(media => ({
        type: 'image',
        src: media.contentPath,
    }));

    return (
        <Lightbox
            open
            close={onClose}
            index={displayedIndex}
            slides={slides}
            carousel={{finite: true, padding: 0}}
            controller={{closeOnBackdropClick: true}}
            plugins={[Slideshow]}
            slideshow={{autoplay: slideshow}}
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
