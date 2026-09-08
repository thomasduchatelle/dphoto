import {Box, Typography} from '@mui/material';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import {Media, MediaType, MediaWithinADay} from '@/domains/catalog/language';
import {Album} from '@/domains/catalog/language';
import {basePath} from '@/libs/requests/basepath';
import Link from '@/components/Link';
import {AlbumCard} from '../../../../../_components/AlbumCard';

export interface AlbumMediaGridProps {
    medias: MediaWithinADay[];
    nextAlbum?: Album;
    previousAlbum?: Album;
}

export function AlbumMediaGrid({medias, nextAlbum, previousAlbum}: AlbumMediaGridProps) {
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

            {/* Next/Previous navigation cards — mobile only (xs/sm/md) */}
            {(nextAlbum || previousAlbum) && (
                <Box
                    sx={{
                        display: {xs: 'flex', lg: 'none'},
                        flexDirection: 'column',
                        gap: 1.5,
                        px: 1.5,
                        pb: 10,
                        pt: 2,
                        borderTop: '1px solid rgba(255,255,255,0.07)',
                    }}
                >
                    {nextAlbum && <NeighbourAlbumLink album={nextAlbum}/>}
                    {previousAlbum && <NeighbourAlbumLink album={previousAlbum}/>}
                </Box>
            )}
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
            <img
                src={`${media.contentPath.startsWith(basePath) ? '' : basePath}${media.contentPath}?w=360`}
                alt=""
                style={{width: '100%', height: '100%', objectFit: 'cover', display: 'block'}}
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

export function NeighbourAlbumLink({album}: {album: Album}) {
    return (
        <Box
            component={Link}
            href={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
            prefetch={false}
            sx={{
                display: 'flex',
                alignItems: 'center',
                gap: 1,
                textDecoration: 'none',
                maxWidth: 340,
            }}
        >
            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumCard album={album} onShare={() => {}} compact/>
            </Box>
            <ChevronRightIcon sx={{color: 'rgba(255,255,255,0.35)', fontSize: 22, flexShrink: 0}}/>
        </Box>
    );
}
