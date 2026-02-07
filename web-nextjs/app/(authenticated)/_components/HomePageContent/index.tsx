'use client';

import {Album} from '@/domains/catalog/language/catalog-state';
import {Box, Button, CircularProgress, Typography} from '@mui/material';
import Link from '@/components/Link';

export interface HomePageContentProps {
    albums: Album[];
    isLoading: boolean;
    error?: Error;
    // TODO Pass the `onRetry` callback when the parent will also be a client component
    // onRetry: () => void;
}

export function HomePageContent({albums, isLoading, error}: HomePageContentProps) {
    if (isLoading) {
        return (
            <Box
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '50vh',
                    gap: 2,
                }}
            >
                <CircularProgress/>
                <Typography>Loading albums...</Typography>
            </Box>
        );
    }

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
        <Box sx={{p: 3}}>
            <Typography variant="h4" sx={{mb: 3}}>
                Albums
            </Typography>
            <Box sx={{display: 'flex', flexDirection: 'column', gap: 1}}>
                {albums.map((album) => (
                    <Box key={`${album.albumId.owner}-${album.albumId.folderName}`} sx={{marginBottom: 1}}>
                        <Link
                            href={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
                            prefetch={false}
                        >
                            <Typography>{album.name}</Typography>
                        </Link>
                    </Box>
                ))}
            </Box>
        </Box>
    );
}
