'use client';

import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {Box, Button, Typography} from '@mui/material';
import {AlbumGrid} from '../AlbumGrid';

export interface HomePageContentProps {
    albums: Album[];
    error?: Error;
    // TODO Pass the `onRetry` callback when the parent will also be a client component
    // onRetry: () => void;
}

export function HomePageContent({albums, error}: HomePageContentProps) {
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

    if (albums.length === 0) {
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
        <AlbumGrid albums={albums} onShare={(id: AlbumId) => console.log('onShare', id)}/>
    );
}
