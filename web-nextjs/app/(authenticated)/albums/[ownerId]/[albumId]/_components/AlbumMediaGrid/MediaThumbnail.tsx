import {Box} from '@mui/material';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import {Media, MediaType} from '@/domains/catalog/language';
import {withBasePath} from '@/libs/requests/media-url';

export interface MediaThumbnailProps {
    media: Media;
}

export function MediaThumbnail({media}: MediaThumbnailProps) {
    const src = media.type === MediaType.VIDEO
        ? withBasePath('/video-placeholder.png')
        : media.thumbnailUrl;

    return (
        <Box
            sx={{
                aspectRatio: '1',
                position: 'relative',
                overflow: 'hidden',
                bgcolor: 'rgba(255,255,255,0.05)',
            }}
        >
            <Box
                component="img"
                src={src}
                alt=""
                loading="lazy"
                sx={{width: '100%', height: '100%', objectFit: 'cover', display: 'block'}}
            />
            {media.type === MediaType.VIDEO && (
                <Box
                    sx={{
                        position: 'absolute',
                        inset: 0,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                    }}
                >
                    <PlayCircleOutlineIcon sx={{color: 'rgba(255,255,255,0.9)', fontSize: {xs: 32, md: 40}}}/>
                </Box>
            )}
        </Box>
    );
}
