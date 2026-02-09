'use client';

import {Box} from '@mui/material';
import {ReactNode} from 'react';

export interface AlbumGridProps {
    children: ReactNode;
}

export const AlbumGrid = ({children}: AlbumGridProps) => {
    return (
        <Box
            component="section"
            aria-label="Album list"
            sx={{
                display: 'grid',
                gridTemplateColumns: {
                    xs: '1fr',
                    sm: 'repeat(2, 1fr)',
                    md: 'repeat(3, 1fr)',
                    lg: 'repeat(4, 1fr)',
                },
                gap: 4,
                width: '100%',
                maxWidth: 1920,
                mx: 'auto',
            }}
        >
            {children}
        </Box>
    );
};
