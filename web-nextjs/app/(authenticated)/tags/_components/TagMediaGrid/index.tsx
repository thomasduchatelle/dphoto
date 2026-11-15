import {Box, Typography} from '@mui/material';
import {TagMediaWithinADay} from '@/domains/catalog/tags';
import {MediaThumbnail} from '@/components/MediaThumbnail';
import {toLocaleDateWithDay} from '@/libs/dates';

export interface TagMediaGridProps {
    medias: TagMediaWithinADay[];
}

export function TagMediaGrid({medias}: TagMediaGridProps) {
    return (
        <Box sx={(theme) => ({
            borderTop: '1px solid rgba(255,255,255,0.07)',
            pt: theme.spacing(2),
        })}>
            {medias.map(({day, medias: dayMedias}) => (
                <Box key={day.toISOString()} sx={{mb: 4}}>
                    <Typography variant="h2" sx={{mb: 1.5, px: {xs: 1, sm: 0}}}>
                        {toLocaleDateWithDay(day)}
                    </Typography>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: {xs: 'repeat(2, 1fr)', sm: 'repeat(3, 1fr)', md: 'repeat(4, 1fr)', lg: 'repeat(5, 1fr)'},
                            gap: '2px',
                        }}
                    >
                        {dayMedias.map((media) => (
                            <MediaThumbnail key={`${media.albumId.owner}/${media.albumId.folderName}/${media.id}`} albumId={media.albumId} media={media}/>
                        ))}
                    </Box>
                </Box>
            ))}
        </Box>
    );
}
