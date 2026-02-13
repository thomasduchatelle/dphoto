'use client';

import {Box, Typography, Badge} from '@mui/material';
import ShareIcon from '@mui/icons-material/Share';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {SharedByIndicator} from '@/components/shared/SharedByIndicator';

export interface AlbumCardProps {
    album: Album;
    onClick: (albumId: AlbumId) => void;
}

// Helper to get temperature color based on relativeTemperature (0-1 scale)
const getTemperatureColor = (relativeTemp: number): string => {
    if (relativeTemp >= 0.8) return '#ff6b6b'; // Hot - red
    if (relativeTemp >= 0.5) return '#ffa94d'; // Warm - orange
    if (relativeTemp >= 0.3) return '#74c0fc'; // Cool - light blue
    return '#a5d8ff'; // Cold - pale blue
};

const formatDateRange = (start: Date, end: Date): string => {
    const formatDate = (date: Date): string => {
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        }).toUpperCase();
    };
    return `${formatDate(start)} – ${formatDate(end)}`;
};

export const AlbumCard = ({album, onClick}: AlbumCardProps) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };
    const tempColor = getTemperatureColor(album.relativeTemperature);

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            aria-label={`Album: ${album.name}, ${album.totalCount} photos`}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                '&:hover': {
                    transform: 'translateY(-6px)',
                    boxShadow: '0 12px 40px rgba(24, 89, 134, 0.4)',
                    '& .album-photo': {
                        filter: 'brightness(0.85)',
                    },
                    '& .compact-bar': {
                        opacity: 0,
                    },
                    '& .expanded-overlay': {
                        opacity: 1,
                    },
                },
                '&:focus': {
                    outline: '2px solid',
                    outlineColor: 'primary.main',
                    outlineOffset: 2,
                },
            }}
        >
            {/* 2x2 Photo Grid */}
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box
                        key={i}
                        className="album-photo"
                        sx={{
                            aspectRatio: '1',
                            background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                            transition: 'filter 0.3s ease',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            fontSize: 32,
                            opacity: 0.15,
                        }}
                    >
                        🖼️
                    </Box>
                ))}
            </Box>

            {/* Compact bar below photos */}
            <Box
                className="compact-bar"
                sx={{
                    position: 'relative',
                    background: 'rgba(10, 21, 32, 0.95)',
                    padding: '12px 16px',
                    transition: 'opacity 0.3s ease',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    '&::before': {
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        right: 0,
                        height: '3px',
                        background: `linear-gradient(to right, ${tempColor}, transparent)`,
                    },
                }}
            >
                <Typography
                    sx={{
                        fontFamily: 'Georgia, serif',
                        fontSize: 18,
                        fontWeight: 300,
                        color: '#ffffff',
                        lineHeight: 1.2,
                        flex: 1,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        mr: 2,
                    }}
                >
                    {album.name}
                </Typography>
                <Box sx={{display: 'flex', alignItems: 'center', gap: 1, flexShrink: 0}}>
                    <Typography
                        sx={{
                            fontFamily: 'Courier New, monospace',
                            fontSize: 13,
                            color: tempColor,
                            fontWeight: 500,
                        }}
                    >
                        {album.totalCount}
                    </Typography>
                    {album.sharedWith.length > 0 && (
                        <Badge
                            badgeContent={album.sharedWith.length}
                            color="primary"
                            sx={{
                                '& .MuiBadge-badge': {
                                    fontSize: 9,
                                    height: 14,
                                    minWidth: 14,
                                    padding: '0 4px',
                                    backgroundColor: tempColor,
                                    color: '#ffffff',
                                    fontWeight: 600,
                                },
                            }}
                        >
                            <ShareIcon sx={{fontSize: 16, color: 'rgba(255, 255, 255, 0.7)'}} />
                        </Badge>
                    )}
                </Box>
            </Box>

            {/* Expanded overlay on hover */}
            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background:
                        'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography
                    sx={{
                        fontFamily: 'Georgia, serif',
                        fontSize: 22,
                        fontWeight: 300,
                        mb: 1.5,
                        color: '#ffffff',
                        lineHeight: 1.3,
                    }}
                >
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{
                            fontFamily: 'Courier New, monospace',
                            fontSize: 13,
                            color: 'rgba(255, 255, 255, 0.85)',
                            fontWeight: 300,
                            letterSpacing: '0.05em',
                        }}
                    >
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography
                        sx={{
                            fontFamily: 'Courier New, monospace',
                            fontSize: 13,
                            color: '#ffffff',
                            fontWeight: 400,
                        }}
                    >
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1,
                            pt: 1,
                            borderTop: '1px solid rgba(255, 255, 255, 0.3)',
                        }}
                    >
                        <Typography
                            sx={{
                                fontSize: 11,
                                color: 'rgba(255, 255, 255, 0.8)',
                                textTransform: 'uppercase',
                                letterSpacing: '0.08em',
                            }}
                        >
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map((s) => s.user)} />
                    </Box>
                )}

                {album.ownedBy && (
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1,
                            pt: 1,
                            borderTop: '1px solid rgba(255, 255, 255, 0.3)',
                            mt: album.sharedWith.length > 0 ? 1 : 0,
                        }}
                    >
                        <Typography
                            sx={{
                                fontSize: 11,
                                color: 'rgba(255, 255, 255, 0.8)',
                                textTransform: 'uppercase',
                                letterSpacing: '0.08em',
                            }}
                        >
                            Shared by
                        </Typography>
                        <SharedByIndicator users={album.ownedBy.users} />
                    </Box>
                )}
            </Box>
        </Box>
    );
};
