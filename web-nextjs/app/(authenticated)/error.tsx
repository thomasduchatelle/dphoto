'use client';

import {Box, Button} from '@mui/material';
import {ErrorDisplay} from '@/components/shared/ErrorDisplay';
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
            <Box>
                <ErrorDisplay
                    error={{
                        message: error.message || 'An unexpected error occurred while loading your albums',
                        details: error.stack,
                    }}
                    onRetry={reset}
                />
                <Box sx={{display: 'flex', justifyContent: 'center', mt: 2}}>
                    <Button component={Link} href="/" variant="outlined" prefetch={false}>
                        Return to Albums
                    </Button>
                </Box>
            </Box>
        </Box>
    );
}
