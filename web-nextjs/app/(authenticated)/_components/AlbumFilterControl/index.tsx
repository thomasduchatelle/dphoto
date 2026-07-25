'use client';

import {Box, ToggleButton, ToggleButtonGroup, Typography} from '@mui/material';
import {AlbumFilterEntry} from '@/domains/catalog/language/catalog-state';

export interface AlbumFilterControlProps {
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterEntry) => void;
}

export function AlbumFilterControl({filterOptions, activeFilter, onFilterChange}: AlbumFilterControlProps) {
    const handleChange = (_event: React.MouseEvent<HTMLElement>, newValue: string | null) => {
        if (newValue === null) return;
        const selectedFilter = filterOptions.find(option => option.name === newValue);
        if (selectedFilter) {
            onFilterChange(selectedFilter);
        }
    };

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: {xs: 'column', sm: 'row'},
                alignItems: {xs: 'stretch', sm: 'center'},
                gap: 1,
                marginBottom: 4,
            }}
        >
            <Typography
                variant="body1"
                color="text.secondary"
                sx={{
                    flexShrink: 0,
                    minWidth: 'fit-content',
                }}
            >
                Filter by owner:
            </Typography>
            <ToggleButtonGroup
                value={activeFilter.name}
                exclusive
                onChange={handleChange}
                aria-label="Album filter by owner"
                sx={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    gap: 0.5,
                    '& .MuiToggleButton-root': {
                        minWidth: '120px',
                        minHeight: '44px',
                        textTransform: 'none',
                        color: 'text.primary',
                        borderColor: 'rgba(255, 255, 255, 0.2)',
                        '&.Mui-selected': {
                            backgroundColor: 'primary.main',
                            color: '#ffffff',
                            '&:hover': {
                                backgroundColor: 'primary.dark',
                            },
                        },
                        '&:focus': {
                            outline: '2px solid',
                            outlineColor: 'primary.main',
                            outlineOffset: '2px',
                        },
                    },
                }}
            >
                {filterOptions.map(option => (
                    <ToggleButton
                        key={option.name}
                        value={option.name}
                        aria-pressed={option.name === activeFilter.name}
                    >
                        {option.name}
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        </Box>
    );
}
