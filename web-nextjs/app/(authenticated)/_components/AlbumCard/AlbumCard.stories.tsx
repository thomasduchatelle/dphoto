import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumCard} from './index';
import {Box, Typography} from '@mui/material';
import {SharedByIndicator} from '@/components/shared/SharedByIndicator';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const sampleAlbums: Album[] = [
    {
        albumId: createAlbumId('owner-1', 'santa-cruz-2026'),
        name: 'Santa Cruz Beach Trip',
        start: new Date('2026-07-15'),
        end: new Date('2026-07-22'),
        totalCount: 47,
        temperature: 6.7,
        relativeTemperature: 0.45,
        sharedWith: [
            {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
            {user: {name: 'Bob Smith', email: 'bob@example.com'}},
        ],
    },
    {
        albumId: createAlbumId('owner-1', 'iceland-2025'),
        name: 'Iceland Adventure',
        start: new Date('2025-08-01'),
        end: new Date('2025-08-14'),
        totalCount: 203,
        temperature: 14.5,
        relativeTemperature: 0.85,
        sharedWith: [],
    },
    {
        albumId: createAlbumId('owner-1', 'paris-2025'),
        name: 'Paris Weekend',
        start: new Date('2025-03-12'),
        end: new Date('2025-03-15'),
        totalCount: 89,
        temperature: 22.3,
        relativeTemperature: 1.0,
        sharedWith: [
            {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
            {user: {name: 'Bob Smith', email: 'bob@example.com'}},
        ],
    },
    {
        albumId: createAlbumId('owner-1', 'sunset-hike'),
        name: 'Mountain Sunset Hike',
        start: new Date('2025-06-08'),
        end: new Date('2025-06-09'),
        totalCount: 24,
        temperature: 2.1,
        relativeTemperature: 0.15,
        sharedWith: [],
    },
];

const formatDate = (date: Date): string => {
    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
    }).toUpperCase();
};

const formatDateRange = (start: Date, end: Date): string => {
    return `${formatDate(start)} – ${formatDate(end)}`;
};

const formatMonth = (date: Date): string => {
    return date.toLocaleDateString('en-US', {month: 'short', year: 'numeric'}).toUpperCase();
};

const meta = {
    title: 'Components/AlbumCard',
    component: AlbumCard,
    parameters: {
        layout: 'fullscreen',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof AlbumCard>;

export default meta;
type Story = StoryObj<typeof meta>;

// ============================================================================
// DESIGN EXPLORATIONS - V1-based variations with photo count, sharing, and temperature
// ============================================================================

// Helper to get temperature color based on relativeTemperature (0-1 scale)
const getTemperatureColor = (relativeTemp: number): string => {
    if (relativeTemp >= 0.8) return '#ff6b6b'; // Hot - red
    if (relativeTemp >= 0.5) return '#ffa94d'; // Warm - orange
    if (relativeTemp >= 0.3) return '#74c0fc'; // Cool - light blue
    return '#a5d8ff'; // Cold - pale blue
};

// V1A: Photo count next to title, sharing icon on right, temperature as border
const AlbumCardV1A = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
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
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                border: `2px solid ${tempColor}`,
                '&:hover': {
                    transform: 'translateY(-6px)',
                    boxShadow: `0 12px 40px ${tempColor}40`,
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
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                <Box sx={{flex: 1, overflow: 'hidden', display: 'flex', alignItems: 'baseline', gap: 1.5}}>
                    <Typography sx={{
                        fontFamily: 'Georgia, serif',
                        fontSize: 18,
                        fontWeight: 300,
                        color: '#ffffff',
                        lineHeight: 1.2,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                    }}>
                        {album.name}
                    </Typography>
                    <Typography sx={{
                        fontFamily: 'Courier New, monospace',
                        fontSize: 12,
                        color: 'rgba(255, 255, 255, 0.6)',
                        flexShrink: 0
                    }}>
                        {album.totalCount}
                    </Typography>
                </Box>
                {album.sharedWith.length > 0 && (
                    <Box sx={{ml: 1, flexShrink: 0, fontSize: 16}}>
                        👥
                    </Box>
                )}
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// V1B: Photo count on right, sharing as small badge, temperature as left accent bar
const AlbumCardV1B = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
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
            {/* Temperature accent bar */}
            <Box sx={{
                position: 'absolute',
                left: 0,
                top: 0,
                bottom: 0,
                width: 4,
                background: tempColor,
                zIndex: 2,
                transition: 'width 0.3s ease',
                '&:hover': {
                    width: 6
                }
            }}/>

            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            {/* Sharing badge overlay on photos */}
            {album.sharedWith.length > 0 && (
                <Box sx={{
                    position: 'absolute',
                    top: 8,
                    right: 8,
                    background: 'rgba(24, 89, 134, 0.95)',
                    borderRadius: '12px',
                    padding: '4px 10px',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 0.5,
                    fontSize: 11,
                    backdropFilter: 'blur(4px)'
                }}>
                    <span>👥</span>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 10, color: '#ffffff'}}>
                        {album.sharedWith.length}
                    </Typography>
                </Box>
            )}

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    flex: 1,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    mr: 2
                }}>
                    {album.name}
                </Typography>
                <Typography sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 13,
                    color: tempColor,
                    fontWeight: 500,
                    flexShrink: 0
                }}>
                    {album.totalCount}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// V1C: Photo count inline, sharing icon + temp dot in compact bar
const AlbumCardV1C = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
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
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                gap: 1.5
            }}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    flex: 1,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                }}>
                    {album.name}
                </Typography>
                <Typography sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 12,
                    color: 'rgba(255, 255, 255, 0.6)',
                    flexShrink: 0
                }}>
                    {album.totalCount}
                </Typography>
                {album.sharedWith.length > 0 && (
                    <Box sx={{fontSize: 14, flexShrink: 0, lineHeight: 1}}>
                        👥
                    </Box>
                )}
                <Box sx={{
                    width: 8,
                    height: 8,
                    borderRadius: '50%',
                    background: tempColor,
                    flexShrink: 0,
                    boxShadow: `0 0 8px ${tempColor}80`
                }}/>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// V1D: Photo count on right, no sharing indicator (only in expanded), temperature as top gradient
const AlbumCardV1D = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
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
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                overflow: 'hidden',
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
            {/* Temperature gradient at top */}
            <Box sx={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                height: 3,
                background: `linear-gradient(to right, ${tempColor}, transparent)`,
                zIndex: 2
            }}/>

            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    flex: 1,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    mr: 2
                }}>
                    {album.name}
                </Typography>
                <Typography sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 13,
                    color: 'rgba(255, 255, 255, 0.7)',
                    flexShrink: 0
                }}>
                    {album.totalCount}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// V1E: Photo count colored by temp, sharing as avatars peek, temp as bottom border
const AlbumCardV1E = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
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
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                borderBottom: `3px solid ${tempColor}`,
                '&:hover': {
                    transform: 'translateY(-6px)',
                    boxShadow: `0 12px 40px ${tempColor}40`,
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
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                <Box sx={{flex: 1, overflow: 'hidden', display: 'flex', alignItems: 'center', gap: 1.5}}>
                    <Typography sx={{
                        fontFamily: 'Georgia, serif',
                        fontSize: 18,
                        fontWeight: 300,
                        color: '#ffffff',
                        lineHeight: 1.2,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                    }}>
                        {album.name}
                    </Typography>
                    {album.sharedWith.length > 0 && (
                        <Box sx={{display: 'flex', alignItems: 'center', ml: 'auto', mr: 1}}>
                            {album.sharedWith.slice(0, 2).map((share, idx) => (
                                <Box
                                    key={idx}
                                    sx={{
                                        width: 20,
                                        height: 20,
                                        borderRadius: '50%',
                                        background: share.user.picture ? `url(${share.user.picture})` : '#4a9ece',
                                        backgroundSize: 'cover',
                                        border: '2px solid rgba(10, 21, 32, 0.95)',
                                        marginLeft: idx > 0 ? '-8px' : 0
                                    }}
                                />
                            ))}
                        </Box>
                    )}
                </Box>
                <Typography sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 14,
                    color: tempColor,
                    fontWeight: 600,
                    flexShrink: 0
                }}>
                    {album.totalCount}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// OLD VARIANTS BELOW (kept for reference)
// ============================================================================

// Variant 1: Title only in compact, full details on hover
const AlbumCardV1 = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.4s ease',
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
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.4s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            {/* Compact bar below photos */}
            <Box className="compact-bar" sx={{background: 'rgba(10, 21, 32, 0.95)', padding: '12px 16px', transition: 'opacity 0.3s ease'}}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                }}>
                    {album.name}
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
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 1.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5}}>
                    <Typography
                        sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', fontWeight: 300, letterSpacing: '0.05em'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// Variant 2: Title + Date range in compact, adds photo count + sharing on hover
const AlbumCardV2 = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.4s ease',
                '&:hover': {
                    transform: 'translateY(-6px)',
                    boxShadow: '0 12px 40px rgba(24, 89, 134, 0.4)',
                    '& .album-photo': {
                        filter: 'brightness(0.85)',
                    },
                    '& .compact-bar': {
                        background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98), rgba(24, 89, 134, 0.95))',
                        paddingTop: '16px',
                        paddingBottom: '16px',
                    },
                    '& .extra-info': {
                        opacity: 1,
                        maxHeight: '100px',
                        marginTop: '12px',
                    },
                },
            }}
        >
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.4s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{background: 'rgba(10, 21, 32, 0.95)', padding: '12px 16px', transition: 'all 0.4s ease'}}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.3,
                    mb: 0.5,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                }}>
                    {album.name}
                </Typography>
                <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', letterSpacing: '0.05em'}}>
                    {formatDateRange(album.start, album.end)}
                </Typography>

                <Box className="extra-info" sx={{opacity: 0, maxHeight: 0, overflow: 'hidden', transition: 'all 0.4s ease'}}>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 12, color: '#ffffff', fontWeight: 400, mb: 1}}>
                        {album.totalCount} photos
                    </Typography>
                    {album.sharedWith.length > 0 && (
                        <Box sx={{display: 'flex', alignItems: 'center', gap: 1}}>
                            <Typography sx={{fontSize: 10, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase'}}>
                                Shared with
                            </Typography>
                            <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                        </Box>
                    )}
                </Box>
            </Box>
        </Box>
    );
};

// Variant 3: Title + Start date in compact, full date range + details on hover
const AlbumCardV3 = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.4s ease',
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
            }}
        >
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.4s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{
                background: 'rgba(10, 21, 32, 0.95)',
                padding: '12px 16px',
                transition: 'opacity 0.3s ease',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center'
            }}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    flex: 1,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    mr: 2
                }}>
                    {album.name}
                </Typography>
                <Typography
                    sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', letterSpacing: '0.05em', flexShrink: 0}}>
                    {formatMonth(album.start)}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 0.5, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: 'rgba(255, 255, 255, 0.85)', mb: 1.5, letterSpacing: '0.05em'}}>
                    {formatDateRange(album.start, album.end)}
                </Typography>

                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                    {album.sharedWith.length > 0 && (
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    )}
                </Box>
            </Box>
        </Box>
    );
};

// Variant 4: Title + Photo count in compact, slides over photos with full info
const AlbumCardV4 = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.4s ease',
                '&:hover': {
                    transform: 'translateY(-6px)',
                    boxShadow: '0 12px 40px rgba(24, 89, 134, 0.4)',
                    '& .album-photo': {
                        filter: 'brightness(0.75)',
                    },
                    '& .expanded-overlay': {
                        transform: 'translateY(0)',
                    },
                },
            }}
        >
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.4s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box sx={{background: 'rgba(10, 21, 32, 0.95)', padding: '12px 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                <Typography sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 18,
                    fontWeight: 300,
                    color: '#ffffff',
                    lineHeight: 1.2,
                    flex: 1,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    mr: 2
                }}>
                    {album.name}
                </Typography>
                <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 12, color: '#6ab9de', fontWeight: 400, flexShrink: 0}}>
                    {album.totalCount}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.85) 50%, rgba(24, 89, 134, 0.7) 100%)',
                    padding: '24px 20px 20px',
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'flex-end',
                    transform: 'translateY(100%)',
                    transition: 'transform 0.4s ease',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 24, fontWeight: 300, mb: 2, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'grid', gridTemplateColumns: '1fr auto', gap: 2, mb: 2}}>
                    <Box>
                        <Typography
                            sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', mb: 0.5}}>
                            Date Range
                        </Typography>
                        <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff'}}>
                            {formatDateRange(album.start, album.end)}
                        </Typography>
                    </Box>
                    <Box sx={{textAlign: 'right'}}>
                        <Typography
                            sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', mb: 0.5}}>
                            Photos
                        </Typography>
                        <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                            {album.totalCount}
                        </Typography>
                    </Box>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1, pt: 1.5, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

// Variant 5: Inline compact with all info, expands over with detailed layout
const AlbumCardV5 = ({album, onClick}: { album: Album; onClick: (albumId: AlbumId) => void }) => {
    const handleClick = () => onClick(album.albumId);
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleClick();
    };

    return (
        <Box
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            sx={{
                position: 'relative',
                cursor: 'pointer',
                transition: 'all 0.4s ease',
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
            }}
        >
            <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 1}}>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i} className="album-photo" sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.4s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 32,
                        opacity: 0.15
                    }}>
                        🖼️
                    </Box>
                ))}
            </Box>

            <Box className="compact-bar" sx={{background: 'rgba(10, 21, 32, 0.95)', padding: '10px 16px', transition: 'opacity 0.3s ease'}}>
                <Box sx={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 0.3}}>
                    <Typography sx={{
                        fontFamily: 'Georgia, serif',
                        fontSize: 16,
                        fontWeight: 300,
                        color: '#ffffff',
                        lineHeight: 1.2,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        flex: 1,
                        mr: 1
                    }}>
                        {album.name}
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: '#6ab9de', flexShrink: 0}}>
                        {album.totalCount}
                    </Typography>
                </Box>
                <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 10, color: 'rgba(255, 255, 255, 0.6)', letterSpacing: '0.05em'}}>
                    {formatMonth(album.start)}
                </Typography>
            </Box>

            <Box
                className="expanded-overlay"
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: 'linear-gradient(to top, rgba(24, 89, 134, 0.98) 0%, rgba(24, 89, 134, 0.75) 70%, transparent 100%)',
                    padding: '48px 20px 20px',
                    opacity: 0,
                    transition: 'opacity 0.3s ease',
                    pointerEvents: 'none',
                }}
            >
                <Typography sx={{fontFamily: 'Georgia, serif', fontSize: 22, fontWeight: 300, mb: 2, color: '#ffffff', lineHeight: 1.3}}>
                    {album.name}
                </Typography>

                <Box sx={{display: 'grid', gridTemplateColumns: 'auto 1fr', gap: 1.5, mb: 2}}>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase'}}>
                        Period
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff'}}>
                        {formatDateRange(album.start, album.end)}
                    </Typography>

                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 11, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase'}}>
                        Count
                    </Typography>
                    <Typography sx={{fontFamily: 'Courier New, monospace', fontSize: 13, color: '#ffffff', fontWeight: 400}}>
                        {album.totalCount} photos
                    </Typography>
                </Box>

                {album.sharedWith.length > 0 && (
                    <Box sx={{display: 'flex', alignItems: 'center', gap: 1.5, pt: 1.5, borderTop: '1px solid rgba(255, 255, 255, 0.3)'}}>
                        <Typography sx={{fontSize: 11, color: 'rgba(255, 255, 255, 0.8)', textTransform: 'uppercase', letterSpacing: '0.08em'}}>
                            Shared with
                        </Typography>
                        <SharedByIndicator users={album.sharedWith.map(s => s.user)}/>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

export const DesignComparison: Story = {
    render: () => (
        <AppBackground>
            <Box sx={{p: 6}}>
                <Typography variant="h1" sx={{
                    mb: 1,
                    fontSize: 28,
                    fontWeight: 300,
                    color: '#4a9ece',
                    textAlign: 'center',
                    letterSpacing: '0.08em',
                    textTransform: 'uppercase'
                }}>
                    AlbumCard V1-Based Variations
                </Typography>
                <Typography sx={{mb: 2, fontSize: 14, color: 'rgba(255, 255, 255, 0.6)', textAlign: 'center'}}>
                    Title + Photo Count in compact · Hover to expand · Different sharing & temperature indicators
                </Typography>
                <Typography sx={{mb: 6, fontSize: 12, color: 'rgba(255, 255, 255, 0.4)', textAlign: 'center', fontStyle: 'italic'}}>
                    Temperature: Cold (pale blue) → Cool (light blue) → Warm (orange) → Hot (red)
                </Typography>

                <Box sx={{display: 'flex', flexDirection: 'column', gap: 8, maxWidth: 1600, mx: 'auto'}}>
                    {/* V1A */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1A: Photo count inline · Sharing icon · Temperature border
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1A key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Count next to title, sharing icon on right (only if shared), temperature as colored border
                        </Typography>
                    </Box>

                    {/* V1B */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1B: Photo count colored · Sharing badge on photos · Temperature accent bar
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1B key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Count on right (temp colored), sharing badge overlay on photos, temperature as left accent bar
                        </Typography>
                    </Box>

                    {/* V1C */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1C: Compact indicators · All in bar · Temperature dot
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1C key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Everything in compact bar: title, count, sharing icon, temperature dot (glowing)
                        </Typography>
                    </Box>

                    {/* V1D */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1D: Minimal compact · Temperature gradient top · Sharing only on hover
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1D key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Clean compact bar (title + count), temperature as top gradient, sharing revealed on hover only
                        </Typography>
                    </Box>

                    {/* V1E */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1E: Avatar peek · Temperature-colored count · Bottom border
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1E key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Shared user avatars peek in bar, count colored by temperature, temperature as bottom border
                        </Typography>
                    </Box>
                </Box>
            </Box>
        </AppBackground>
    ),
};
