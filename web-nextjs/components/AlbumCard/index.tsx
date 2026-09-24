'use client';

import {Box} from '@mui/material';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {AlbumThumbnails} from './AlbumThumbnails';
import {CompactBar} from './CompactBar';
import {ExpandedOverlay} from './ExpandedOverlay';

export interface AlbumCardProps {
    album: Album;
    onShare?: (albumId: AlbumId) => void;
    compact?: boolean;
}

const getTemperatureColor = (relativeTemp: number): string => {
    if (relativeTemp >= 0.75) return '#ff6b6b';
    if (relativeTemp >= 0.5) return '#ffa94d';
    if (relativeTemp >= 0.25) return '#74c0fc';
    return '#a5d8ff';
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

export const AlbumCard = ({album, onShare, compact = false}: AlbumCardProps) => {
    const temperatureColor = getTemperatureColor(album.relativeTemperature);
    const dateRange = formatDateRange(album.start, album.end);
    const onShareClick = onShare
        ? (e: React.MouseEvent) => {
            e.preventDefault();
            e.stopPropagation();
            onShare(album.albumId);
        }
        : undefined;

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
                            opacity: compact ? 100 : 0,
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
            <AlbumThumbnails album={album} compact={compact}/>
            <CompactBar album={album} temperatureColor={temperatureColor} dateRange={dateRange} compact={compact}/>
            {!compact && <ExpandedOverlay album={album} temperatureColor={temperatureColor} dateRange={dateRange} onShareClick={onShareClick}/>}
        </Box>
    );
};
