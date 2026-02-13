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
        relativeTemperature: 0.5,
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
        relativeTemperature: 0.9,
        sharedWith: [
            {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
            {user: {name: 'Bob Smith', email: 'bob@example.com'}},
        ],
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
        albumId: createAlbumId('owner-1', 'christmas-2024'),
        name: 'Christmas 2024',
        start: new Date('2024-12-23'),
        end: new Date('2024-12-26'),
        totalCount: 156,
        temperature: 39.0,
        relativeTemperature: 1.0,
        sharedWith: [
            {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
            {user: {name: 'Bob Smith', email: 'bob@example.com'}},
        ],
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
// DESIGN EXPLORATIONS - Compact overlay variations
// ============================================================================

// ============================================================================
// DESIGN EXPLORATIONS - Compact overlay variations
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
                    AlbumCard Compact Overlay Variations
                </Typography>
                <Typography sx={{mb: 6, fontSize: 14, color: 'rgba(255, 255, 255, 0.6)', textAlign: 'center'}}>
                    Using Album domain type · Hover to expand · Compare different information layouts
                </Typography>

                <Box sx={{display: 'flex', flexDirection: 'column', gap: 8, maxWidth: 1600, mx: 'auto'}}>
                    {/* Variant 1 */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V1: Title Only → Full Overlay
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV1 key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Minimal bar with title only, expands to show all details over photos
                        </Typography>
                    </Box>

                    {/* Variant 2 */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V2: Title + Date → Bar Expands
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV2 key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Bar shows title + dates, expands in-place to reveal count + sharing
                        </Typography>
                    </Box>

                    {/* Variant 3 */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V3: Title + Month → Full Overlay
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV3 key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Compact bar with title + month, expands to full date range + inline sharing
                        </Typography>
                    </Box>

                    {/* Variant 4 */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V4: Title + Count → Full Screen Slide
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV4 key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Bar with title + count, slides entire overlay up covering all photos with grid layout
                        </Typography>
                    </Box>

                    {/* Variant 5 */}
                    <Box>
                        <Typography sx={{mb: 2, fontSize: 13, color: 'rgba(255, 255, 255, 0.7)', textTransform: 'uppercase', letterSpacing: '0.1em', pl: 2}}>
                            V5: Inline All Info → Detailed Grid
                        </Typography>
                        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>
                            {sampleAlbums.map(album => (
                                <AlbumCardV5 key={album.albumId.folderName} album={album} onClick={fn()}/>
                            ))}
                        </Box>
                        <Typography sx={{mt: 1.5, fontSize: 11, color: 'rgba(255, 255, 255, 0.4)', fontStyle: 'italic', pl: 2}}>
                            Compact two-line bar with all basics, expands to labeled grid layout with sharing
                        </Typography>
                    </Box>
                </Box>
            </Box>
        </AppBackground>
    ),
};
