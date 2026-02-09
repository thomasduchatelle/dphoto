'use client';

import {Box, CircularProgress} from '@mui/material';

export interface NavigationLoadingIndicatorProps {
    variant?: 'inline' | 'overlay';
}

export const NavigationLoadingIndicator = ({variant = 'inline'}: NavigationLoadingIndicatorProps) => {
    if (variant === 'overlay') {
        return (
            <Box
                sx={{
                    position: 'fixed',
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    bgcolor: 'rgba(0, 0, 0, 0.5)',
                    zIndex: 1300,
                }}
            >
                <CircularProgress
                    size={20}
                    sx={{
                        color: 'primary.main',
                    }}
                />
            </Box>
        );
    }

    return (
        <Box
            sx={{
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            <CircularProgress
                size={20}
                sx={{
                    color: 'primary.main',
                }}
            />
        </Box>
    );
};
