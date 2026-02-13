'use client';

import {Badge, Box, IconButton, Typography} from '@mui/material';
import ShareIcon from '@mui/icons-material/Share';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {SharedByIndicator} from './SharedByIndicator';

export interface AlbumCardProps {
    album: Album;
    onShare: (albumId: AlbumId) => void;
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
        return date.toLocaleDateString('en-GB', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        }).toUpperCase();
    };
    return `${formatDate(start)} → ${formatDate(end)}`;
};

export const AlbumCard = ({album, onShare}: AlbumCardProps) => {
    const handleClickOnShare = (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        onShare(album.albumId);
    };
    const temperatureColor = getTemperatureColor(album.relativeTemperature);

    return (
        <Box
            aria-label={`Album: ${album.name}, ${album.totalCount} photos`}
            sx={(theme) => ({
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                [theme.breakpoints.up('sm')]: {
                    '&:hover': {
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
                },
                '&:focus': {
                    outline: '2px solid',
                    outlineColor: 'primary.main',
                    outlineOffset: 2,
                },
            })}
        >
            {/* Photo grid: 2x2 on desktop, 1x4 on mobile */}
            <Box sx={{display: 'grid', gridTemplateColumns: {xs: 'repeat(4, 1fr)', sm: 'repeat(2, 1fr)'}, gap: 1}}>
                {[0, 1, 2, 3].map((i) => {
                    const thumbnail = album.thumbnails?.[i];
                    return (
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
                                fontSize: !thumbnail ? 32 : 12,
                                opacity: !thumbnail ? 0.15 : 1,
                                overflow: 'hidden',
                                position: 'relative',
                            }}
                        >
                            {thumbnail ? (
                                <img
                                    src={thumbnail}
                                    alt={`${album.name} thumbnail ${i + 1}`}
                                    style={{
                                        width: '100%',
                                        height: '100%',
                                        objectFit: 'cover',
                                    }}
                                />
                            ) : (
                                '🖼️'
                            )}
                        </Box>
                    );
                })}
            </Box>

            {/* Compact bar below photos */}
            <Box
                className="compact-bar"
                sx={{
                    position: 'relative',
                    background: 'rgba(10, 21, 32, 0.95)',
                    padding: '12px 15px 15px 5px',
                    transition: 'opacity 0.3s ease',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    flexWrap: 'wrap',
                    gap: 0.5,
                    '&::before': {
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        right: 0,
                        height: '3px',
                        background: `linear-gradient(to right, ${temperatureColor}, transparent)`,
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
                            color: temperatureColor,
                            fontWeight: 500,
                        }}
                    >
                        {album.totalCount}
                    </Typography>
                    {!album.ownedBy && (
                        <IconButton
                            onClick={handleClickOnShare}
                            sx={{display: {xs: 'inline-flex', sm: 'none'}, p: 0.5}}
                        >
                            <Badge
                                badgeContent={album.sharedWith.length}
                                color="primary"
                                sx={{
                                    '& .MuiBadge-badge': {
                                        fontSize: 9,
                                        height: 14,
                                        minWidth: 14,
                                        padding: '0 4px',
                                        backgroundColor: temperatureColor,
                                        color: '#ffffff',
                                        fontWeight: 600,
                                    },
                                }}
                            >
                                <ShareIcon sx={{fontSize: 16, color: 'rgba(255, 255, 255, 0.7)'}}/>
                            </Badge>
                        </IconButton>
                    )}
                    {album.sharedWith.length > 0 && (
                        <Badge
                            badgeContent={album.sharedWith.length}
                            color="primary"
                            sx={{
                                display: {xs: 'none', sm: 'inline-flex'},
                                '& .MuiBadge-badge': {
                                    fontSize: 9,
                                    height: 14,
                                    minWidth: 14,
                                    padding: '0 4px',
                                    backgroundColor: temperatureColor,
                                    color: '#ffffff',
                                    fontWeight: 600,
                                },
                            }}
                        >
                            <ShareIcon sx={{fontSize: 16, color: 'rgba(255, 255, 255, 0.7)'}}/>
                        </Badge>
                    )}
                    {album.ownedBy && (
                        <SharedByIndicator users={album.ownedBy.users}/>
                    )}
                </Box>
                {/* Date range: only visible on mobile */}
                <Typography
                    sx={{
                        display: {xs: 'block', sm: 'none'},
                        width: '100%',
                        fontFamily: 'Courier New, monospace',
                        fontSize: 11,
                        color: 'rgba(255, 255, 255, 0.65)',
                        fontWeight: 300,
                        letterSpacing: '0.05em',
                    }}
                >
                    {formatDateRange(album.start, album.end)}
                </Typography>
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
                    padding: '48px 15px 15px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                }}
            >

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', mb: 1.5}}>
                    <Typography
                        sx={{
                            fontFamily: 'Georgia, serif',
                            fontSize: 22,
                            fontWeight: 300,
                            color: '#ffffff',
                            lineHeight: 1.3,
                        }}
                    >
                        {album.name}
                    </Typography>
                    <Typography
                        sx={{
                            fontFamily: 'Courier New, monospace',
                            fontSize: 13,
                            color: temperatureColor,
                            fontWeight: 400,
                            lineHeight: 1.3,
                        }}
                    >
                        {album.totalCount} medias
                    </Typography>
                </Box>

                <Box sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    borderTop: `1px solid ${temperatureColor}`,
                    pt: 1,
                    mb: 1.5,
                }}>
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
                </Box>

                {!album.ownedBy && (
                    <Box
                        sx={{
                            display: 'flex',
                            flexWrap: "nowrap",
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            gap: 1,
                            pt: 1,
                        }}
                    >
                        {album.sharedWith.length > 0 && (
                            <Box sx={{display: 'flex', alignItems: 'center', gap: 1}}>
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
                                <SharedByIndicator users={album.sharedWith.map((s) => s.user)}/>
                            </Box>
                        ) || (
                            <Typography
                                sx={{
                                    fontSize: 11,
                                    color: 'rgba(255, 255, 255, 0.6)',
                                    fontStyle: 'italic',
                                }}
                            >
                                Private album
                            </Typography>
                        )}
                        <IconButton
                            onClick={handleClickOnShare}>
                            <Badge
                                badgeContent={album.sharedWith.length}
                                color="primary"
                                sx={{
                                    '& .MuiBadge-badge': {
                                        fontSize: 9,
                                        height: 14,
                                        minWidth: 14,
                                        padding: '0 4px',
                                        backgroundColor: temperatureColor,
                                        color: '#ffffff',
                                        fontWeight: 600,
                                    },
                                }}
                            >
                                <ShareIcon sx={{fontSize: 16, color: 'rgba(255, 255, 255, 0.7)'}}/>
                            </Badge>
                        </IconButton>
                    </Box>
                )}

                {album.ownedBy && (
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'space-between',
                            gap: 1,
                            pt: 1,
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
                        <SharedByIndicator users={album.ownedBy.users}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};
