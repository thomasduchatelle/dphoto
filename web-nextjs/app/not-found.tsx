'use client';

import {Box} from '@mui/material';
import {SearchOff as SearchOffIcon} from '@mui/icons-material';
import {EmptyState} from '@/components/shared/EmptyState';

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
                action={{
                    label: 'Go Home',
                    onClick: () => {
                        window.location.href = '/';
                    },
                }}
            />
        </Box>
    );
}
