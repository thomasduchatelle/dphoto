'use client';

import {Box, IconButton} from '@mui/material';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import DownloadIcon from '@mui/icons-material/Download';
import Link from '@/components/Link';
import {AlbumId, Media, MediaType} from '@/domains/catalog/language';
import {photoUrl} from '@/domains/catalog/navigation/photo-url';
import {toLocaleTimeOfDay} from '@/libs/dates';

export interface MediaThumbnailProps {
    albumId: AlbumId;
    media: Media;
}

export function MediaThumbnail({albumId, media}: MediaThumbnailProps) {
    const downloadUrl = media.contentPath;
    const downloadName = media.contentPath.split('/').pop() || media.id;

    const handleDownload = (e: React.MouseEvent) => {
        // The thumbnail is itself a link; stop navigation and trigger a plain
        // download of the original (unresized) media without nesting anchors.
        e.preventDefault();
        e.stopPropagation();
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = downloadName;
        document.body.appendChild(link);
        link.click();
        link.remove();
    };

    return (
        <Box
            component={Link}
            href={photoUrl(albumId, media.id)}
            prefetch={false}
            sx={{
                aspectRatio: '1',
                position: 'relative',
                overflow: 'hidden',
                bgcolor: 'rgba(255,255,255,0.05)',
                display: 'block',
                '& .MediaThumbnail-overlay': {
                    display: 'none',
                },
                '&:hover .MediaThumbnail-overlay': {
                    display: {lg: 'flex'},
                },
            }}
        >
            <Box
                component="img"
                src={media.thumbnailUrl}
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
            <Box
                className="MediaThumbnail-overlay"
                sx={{
                    position: 'absolute',
                    left: 0,
                    right: 0,
                    bottom: 0,
                    alignItems: 'center',
                    gap: 0.5,
                    pl: 1.5,
                    pr: 0.5,
                    pt: 1,
                    pb: 0.5,
                    color: 'common.white',
                    bgcolor: 'rgba(0,0,0,0.5)',
                }}
            >
                <Box
                    sx={{
                        flex: 1,
                        minWidth: 0,
                        fontSize: '0.6em',
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                    }}
                >
                    {media.source ? `@${media.source}` : ''}
                </Box>
                <Box sx={{fontSize: '0.8em', whiteSpace: 'nowrap'}}>{toLocaleTimeOfDay(media.time)}</Box>
                <Box sx={{flex: 1, display: 'flex', justifyContent: 'flex-end', minWidth: 0}}>
                    <IconButton
                        aria-label="Download original media"
                        title="Download"
                        size="small"
                        onClick={handleDownload}
                        sx={{color: 'rgba(255,255,255,0.7)', p: 0.5}}
                    >
                        <DownloadIcon sx={{fontSize: '1rem'}}/>
                    </IconButton>
                </Box>
            </Box>
        </Box>
    );
}
