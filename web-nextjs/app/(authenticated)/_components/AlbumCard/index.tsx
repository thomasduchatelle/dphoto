'use client';

import {Box, Card, Typography} from '@mui/material';
import {UserAvatar} from '@/components/UserAvatar';
import {SharedByIndicator} from '@/components/shared/SharedByIndicator';

export interface AlbumCardProps {
    album: {
        albumId: string;
        ownerId: string;
        name: string;
        startDate: string;
        endDate: string;
        mediaCount: number;
    };
    owner?: {
        name: string;
        email: string;
        picture?: string;
    };
    sharedWith?: Array<{
        name: string;
        email: string;
        picture?: string;
    }>;
    onClick: (albumId: string, ownerId: string) => void;
}

type Density = 'high' | 'medium' | 'low';

const DENSITY_COLORS: Record<Density, string> = {
    high: '#ff6b6b',
    medium: '#ffd43b',
    low: '#51cf66',
};

const calculateDensity = (startDate: string, endDate: string, mediaCount: number): Density => {
    const days = Math.ceil((new Date(endDate).getTime() - new Date(startDate).getTime()) / (1000 * 60 * 60 * 24));
    const photosPerDay = mediaCount / days;
    if (photosPerDay > 10) return 'high';
    if (photosPerDay >= 3) return 'medium';
    return 'low';
};

const formatDateRange = (startDate: string, endDate: string): string => {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const format = (date: Date) =>
        date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    return `${format(start)} - ${format(end)}`.toUpperCase();
};

export const AlbumCard = ({album, owner, sharedWith, onClick}: AlbumCardProps) => {
    const density = calculateDensity(album.startDate, album.endDate, album.mediaCount);
    const densityColor = DENSITY_COLORS[density];

    const handleClick = () => {
        onClick(album.albumId, album.ownerId);
    };

    const handleKeyDown = (event: React.KeyboardEvent) => {
        if (event.key === 'Enter') {
            handleClick();
        }
    };

    return (
        <Card
            role="button"
            tabIndex={0}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            aria-label={`Album: ${album.name}, ${album.mediaCount} photos, ${formatDateRange(album.startDate, album.endDate)}`}
            sx={{
                bgcolor: 'background.paper',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                borderRadius: 0,
                p: 3,
                cursor: 'pointer',
                transition: 'box-shadow 0.2s',
                '&:hover': {
                    boxShadow: 4,
                },
                '&:focus': {
                    outline: `2px solid`,
                    outlineColor: 'primary.main',
                    outlineOffset: 2,
                },
            }}
        >
            <Typography
                variant="h2"
                sx={{
                    fontFamily: 'Georgia, serif',
                    fontSize: 22,
                    fontWeight: 300,
                    mb: 1,
                }}
            >
                {album.name}
            </Typography>

            <Typography
                sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 13,
                    textTransform: 'uppercase',
                    letterSpacing: '0.1em',
                    color: 'text.secondary',
                    mb: 0.5,
                }}
            >
                {formatDateRange(album.startDate, album.endDate)}
            </Typography>

            <Typography
                sx={{
                    fontFamily: 'Courier New, monospace',
                    fontSize: 13,
                    textTransform: 'uppercase',
                    letterSpacing: '0.1em',
                    color: densityColor,
                    mb: 2,
                }}
            >
                {album.mediaCount} PHOTOS
            </Typography>

            {owner && (
                <Box sx={{display: 'flex', alignItems: 'center', gap: 1, mb: 1}}>
                    <UserAvatar name={owner.name} picture={owner.picture} size="small"/>
                    <Typography variant="body2" color="text.secondary">
                        {owner.name}
                    </Typography>
                </Box>
            )}

            {sharedWith && sharedWith.length > 0 && (
                <Box sx={{mt: 2}}>
                    <SharedByIndicator users={sharedWith}/>
                </Box>
            )}
        </Card>
    );
};
