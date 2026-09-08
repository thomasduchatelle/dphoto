'use client';

import {Box, IconButton, Typography} from '@mui/material';
import {Album} from '@/domains/catalog/language/catalog-state';
import {SharedByIndicator} from './SharedByIndicator';
import {ShareBadge} from './ShareBadge';

export interface SharingSummaryProps {
    album: Album;
    temperatureColor: string;
    onShareClick?: (e: React.MouseEvent) => void;
}

const labelSx = {
    fontSize: 11,
    color: 'rgba(255, 255, 255, 0.8)',
    textTransform: 'uppercase',
    letterSpacing: '0.08em',
};

export const SharingSummary = ({album, temperatureColor, onShareClick}: SharingSummaryProps) => {
    if (album.ownedBy) {
        return (
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
                <Typography sx={labelSx}>
                    Shared by
                </Typography>
                <SharedByIndicator users={album.ownedBy.users}/>
            </Box>
        );
    }

    return (
        <Box
            sx={{
                display: 'flex',
                flexWrap: 'nowrap',
                justifyContent: 'space-between',
                alignItems: 'center',
                gap: 1,
                pt: 1,
            }}
        >
            {album.sharedWith.length > 0 ? (
                <Box sx={{display: 'flex', alignItems: 'center', gap: 1}}>
                    <Typography sx={labelSx}>
                        Shared with
                    </Typography>
                    <SharedByIndicator users={album.sharedWith.map((s) => s.user)}/>
                </Box>
            ) : (
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
            {onShareClick && (
                <IconButton onClick={onShareClick}>
                    <ShareBadge count={album.sharedWith.length} color={temperatureColor}/>
                </IconButton>
            )}
        </Box>
    );
};
