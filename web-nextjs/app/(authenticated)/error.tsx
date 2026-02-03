'use client';

import {Box, Button, Paper, Typography} from '@mui/material';
import {ErrorOutline as ErrorOutlineIcon} from '@mui/icons-material';
import Link from '@/components/Link';

export default function AuthenticatedError({
                                               error,
                                               reset,
                                           }: {
    error: Error & { digest?: string }
    reset: () => void
}) {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                p: 2,
            }}
        >
            <Paper
                elevation={3}
                sx={{
                    p: 4,
                    maxWidth: 500,
                    textAlign: 'center',
                }}
            >
                <ErrorOutlineIcon
                    sx={{
                        fontSize: 64,
                        color: 'error.main',
                        mb: 2,
                    }}
                />
                <Typography variant="h4" gutterBottom>
                    Something went wrong
                </Typography>
                <Typography variant="body1" color="text.secondary" sx={{mb: 3}}>
                    {error.message || 'An unexpected error occurred while loading your albums'}
                </Typography>
                <Box sx={{display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap'}}>
                    <Button
                        variant="contained"
                        onClick={reset}
                        sx={{bgcolor: 'primary.main'}}
                    >
                        Try Again
                    </Button>
                    <Link href="/">
                        <Button variant="outlined">
                            Return to Albums
                        </Button>
                    </Link>
                </Box>
            </Paper>
        </Box>
    );
}
