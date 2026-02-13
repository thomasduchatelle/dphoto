'use client';

import {Box} from '@mui/material';
import {ReactNode} from 'react';

export interface AppBackgroundProps {
    children: ReactNode;
}

export function AppBackground({children}: AppBackgroundProps) {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                flexDirection: 'column',
                background: 'linear-gradient(135deg, #001929 0%, #0a1520 25%, #12242e 50%, #0f1d28 75%, #001929 100%)',
                backgroundAttachment: 'fixed',
            }}
        >
            {children}
        </Box>
    );
}
