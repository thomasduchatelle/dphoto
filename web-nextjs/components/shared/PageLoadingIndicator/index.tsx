'use client';

import {Box, LinearProgress, Typography} from '@mui/material';

export interface PageLoadingIndicatorProps {
    message?: string;
}

export const PageLoadingIndicator = ({message = 'Loading...'}: PageLoadingIndicatorProps) => {
    return (
        <Box role="status" aria-live="polite">
            <LinearProgress
                sx={{
                    position: 'fixed',
                    top: 0,
                    left: 0,
                    right: 0,
                    zIndex: 1200,
                    height: 3,
                    bgcolor: 'transparent',
                    '& .MuiLinearProgress-bar': {
                        bgcolor: 'primary.main',
                    },
                }}
            />
            {message && (
                <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{
                        textAlign: 'center',
                        mt: 2,
                    }}
                >
                    {message}
                </Typography>
            )}
        </Box>
    );
};
