'use client';

import {Box, Button} from '@mui/material';
import {SearchOff as SearchOffIcon} from '@mui/icons-material';
import {EmptyState, emptyStateButtonStyles} from '@/components/shared/EmptyState';

export default function NotFound() {
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
            <EmptyState
                icon={<SearchOffIcon/>}
                title="Page Not Found"
                message="The page you're looking for doesn't exist."
            >
                <Button
                    variant="contained"
                    onClick={() => {
                        window.location.href = '/';
                    }}
                    sx={emptyStateButtonStyles.contained}
                >
                    Go Home
                </Button>
            </EmptyState>
        </Box>
    );
}
