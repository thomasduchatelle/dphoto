'use client';

import {Box, Typography} from '@mui/material';
import {Album} from '@/domains/catalog/language/catalog-state';
import {SharingSummary} from './SharingSummary';

export interface ExpandedOverlayProps {
    album: Album;
    temperatureColor: string;
    dateRange: string;
    onShareClick?: (e: React.MouseEvent) => void;
}

export const ExpandedOverlay = ({album, temperatureColor, dateRange, onShareClick}: ExpandedOverlayProps) => (
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
                {dateRange}
            </Typography>
        </Box>

        <SharingSummary album={album} temperatureColor={temperatureColor} onShareClick={onShareClick}/>
    </Box>
);
