'use client';

import {Album, AlbumFilterEntry, AlbumId} from '@/domains/catalog/language/catalog-state';
import {Box, Button, Typography} from '@mui/material';
import {AlbumGrid} from '../AlbumGrid';
import {AlbumFilterControl} from '../AlbumFilterControl';
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION} from '@/domains/catalog/common/utils';

export interface HomePageContentProps {
    albums: Album[];
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterEntry) => void;
    error?: Error;
}

function isAllAlbumsFilter(filter: AlbumFilterEntry): boolean {
    return albumFilterAreCriterionEqual(filter.criterion, ALL_ALBUMS_FILTER_CRITERION);
}

export function HomePageContent({albums, filterOptions, activeFilter, onFilterChange, error}: HomePageContentProps) {
    if (error) {
        return (
            <Box
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '50vh',
                    gap: 2,
                    p: 3,
                }}
            >
                <Typography color="error" variant="h6">
                    Error loading albums
                </Typography>
                <Typography color="text.secondary">
                    {error.message}
                </Typography>
                <Button variant="contained" onClick={() => window.location.reload()}>
                    Try Again
                </Button>
            </Box>
        );
    }

    if (albums.length === 0 && isAllAlbumsFilter(activeFilter)) {
        return (
            <Box
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '50vh',
                    gap: 2,
                    p: 3,
                }}
            >
                <Typography variant="h6">
                    No albums found
                </Typography>
                <Typography color="text.secondary">
                    Create your first album to get started.
                </Typography>
            </Box>
        );
    }

    return (
        <Box>
            <Box sx={{width: '100%', maxWidth: 1920, mx: 'auto', marginBottom: 4}}>
                <AlbumFilterControl
                    filterOptions={filterOptions}
                    activeFilter={activeFilter}
                    onFilterChange={onFilterChange}
                />
            </Box>
            {albums.length === 0 ? (
                <Typography sx={{marginTop: 4, textAlign: 'center', color: 'text.secondary'}}>
                    No albums match this filter.
                </Typography>
            ) : (
                <AlbumGrid albums={albums} onShare={(id: AlbumId) => console.log('onShare', id)}/>
            )}
        </Box>
    );
}
