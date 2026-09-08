'use client';

import {Box, IconButton, Typography} from '@mui/material';
import {Album} from '@/domains/catalog/language/catalog-state';
import {SharedByIndicator} from './SharedByIndicator';
import {ShareBadge} from './ShareBadge';

export interface CompactBarProps {
    album: Album;
    temperatureColor: string;
    dateRange: string;
    onShareClick?: (e: React.MouseEvent) => void;
}

export const CompactBar = ({album, temperatureColor, dateRange, onShareClick}: CompactBarProps) => (
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
            {!album.ownedBy && onShareClick && (
                <IconButton
                    onClick={onShareClick}
                    sx={{display: {xs: 'inline-flex', sm: 'none'}, p: 0.5}}
                >
                    <ShareBadge count={album.sharedWith.length} color={temperatureColor}/>
                </IconButton>
            )}
            {album.sharedWith.length > 0 && (
                <ShareBadge
                    count={album.sharedWith.length}
                    color={temperatureColor}
                    sx={{display: {xs: 'none', sm: 'inline-flex'}}}
                />
            )}
            {album.ownedBy && (
                <SharedByIndicator users={album.ownedBy.users}/>
            )}
        </Box>
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
            {dateRange}
        </Typography>
    </Box>
);
