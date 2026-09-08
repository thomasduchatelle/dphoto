import {Box, Typography} from '@mui/material';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import Image from 'next/image';
import {Media, MediaType, MediaWithinADay} from '@/domains/catalog/language';

export interface AlbumMediaGridProps {
    medias: MediaWithinADay[];
}

export function AlbumMediaGrid({medias}: AlbumMediaGridProps) {
    return (
        <Box>
            {medias.map(({day, medias: dayMedias}) => (
                <Box key={day.toISOString()} sx={{mb: 4}}>
                    <Typography
                        variant="h2"
                        sx={{mb: 1.5}}
                    >
                        {day.toLocaleDateString(undefined, {weekday: 'long', month: 'long', day: 'numeric', year: 'numeric'})}
                    </Typography>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: {xs: 'repeat(2, 1fr)', sm: 'repeat(3, 1fr)', md: 'repeat(4, 1fr)', lg: 'repeat(5, 1fr)'},
                            gap: '2px',
                        }}
                    >
                        {dayMedias.map((media) => (
                            <MediaThumbnail key={media.id} media={media}/>
                        ))}
                    </Box>
                </Box>
            ))}
        </Box>
    );
}

function MediaThumbnail({media}: {media: Media}) {
    return (
        <Box
            sx={{
                aspectRatio: '1',
                position: 'relative',
                overflow: 'hidden',
                cursor: 'pointer',
                bgcolor: 'rgba(255,255,255,0.05)',
                '&:hover': {opacity: 0.85},
                transition: 'opacity 0.15s',
            }}
        >
            <Image
                src={media.contentPath}
                alt=""
                fill
                sizes="(max-width: 600px) 50vw, (max-width: 960px) 33vw, (max-width: 1280px) 25vw, 20vw"
                style={{objectFit: 'cover'}}
            />
            {media.type === MediaType.VIDEO && (
                <Box
                    sx={{
                        position: 'absolute',
                        inset: 0,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        bgcolor: 'rgba(0,0,0,0.25)',
                    }}
                >
                    <PlayCircleOutlineIcon sx={{color: 'rgba(255,255,255,0.9)', fontSize: {xs: 32, md: 40}}}/>
                </Box>
            )}
        </Box>
    );
}
