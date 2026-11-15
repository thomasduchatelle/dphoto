'use client';

import {Box, IconButton} from '@mui/material';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import DownloadIcon from '@mui/icons-material/Download';
import StarIcon from '@mui/icons-material/Star';
import StarBorderIcon from '@mui/icons-material/StarBorder';
import Link from '@/components/Link';
import {AlbumId, FAVORITE_TAG, Media, MediaId, MediaType} from '@/domains/catalog/language';
import {photoUrl} from '@/domains/catalog/navigation/photo-url';
import {toLocaleTimeOfDay} from '@/libs/dates';

export interface MediaThumbnailProps {
    albumId: AlbumId;
    media: Media;
    isOwner: boolean;
    onToggleFavorite: (mediaId: MediaId, currentTags: string[]) => void;
}

export function MediaThumbnail({albumId, media, isOwner, onToggleFavorite}: MediaThumbnailProps) {
    const downloadUrl = media.contentPath;
    const downloadName = media.contentPath.split('/').pop() || media.id;
    const isFavorite = media.tags.includes(FAVORITE_TAG);

    const handleToggleFavorite = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        onToggleFavorite(media.id, media.tags);
    };

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
                '& .MediaThumbnail-star': {
                    display: isFavorite ? 'flex' : 'none',
                },
                '&:hover .MediaThumbnail-star': {
                    display: 'flex',
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
            {isOwner && (
                <Box
                    className="MediaThumbnail-star"
                    sx={{
                        position: 'absolute',
                        top: 0,
                        right: 0,
                        alignItems: 'flex-start',
                        justifyContent: 'flex-end',
                        p: 0.25,
                    }}
                >
                    <IconButton
                        aria-label={isFavorite ? 'Remove from favourites' : 'Mark as favourite'}
                        aria-pressed={isFavorite}
                        title={isFavorite ? 'Remove from favourites' : 'Mark as favourite'}
                        size="small"
                        onClick={handleToggleFavorite}
                        sx={{
                            color: isFavorite ? 'warning.main' : 'rgba(255,255,255,0.85)',
                            p: 0.5,
                            filter: 'drop-shadow(0 1px 2px rgba(0,0,0,0.6))',
                        }}
                    >
                        {isFavorite
                            ? <StarIcon sx={{fontSize: '1.2rem'}}/>
                            : <StarBorderIcon sx={{fontSize: '1.2rem'}}/>}
                    </IconButton>
                </Box>
            )}
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
