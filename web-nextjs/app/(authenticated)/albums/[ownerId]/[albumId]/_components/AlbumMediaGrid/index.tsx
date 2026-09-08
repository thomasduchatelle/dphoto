import {Box, Typography} from '@mui/material';
import {MediaWithinADay} from '@/domains/catalog/language';
import {MediaThumbnail} from './MediaThumbnail';

export interface AlbumMediaGridProps {
    medias: MediaWithinADay[];
}

export function AlbumMediaGrid({medias}: AlbumMediaGridProps) {
    return (
        <Box>
            {medias.map(({day, medias: dayMedias}) => (
                <Box key={day.toISOString()} sx={{mb: 4}}>
                    <Typography variant="h2" sx={{mb: 1.5, px: {xs: 1, sm: 0}}}>
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
