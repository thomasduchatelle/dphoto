'use client';

import {Box, Button, Typography} from '@mui/material';
import {ReactNode} from 'react';

export interface EmptyStateProps {
    icon?: ReactNode;
    title: string;
    message: string;
    action?: {
        label: string;
        onClick: () => void;
    };
}

export const EmptyState = ({icon, title, message, action}: EmptyStateProps) => {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                textAlign: 'center',
                maxWidth: 400,
                mx: 'auto',
                p: 6,
            }}
        >
            {icon && (
                <Box
                    sx={{
                        fontSize: 48,
                        color: 'text.secondary',
                        mb: 2,
                        '& > svg': {
                            fontSize: 48,
                        },
                    }}
                >
                    {icon}
                </Box>
            )}
            <Typography variant="h5" component="h2" gutterBottom>
                {title}
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{mb: action ? 3 : 0}}>
                {message}
            </Typography>
            {action && (
                <Button variant="contained" onClick={action.onClick} sx={{bgcolor: 'primary.main'}}>
                    {action.label}
                </Button>
            )}
        </Box>
    );
};
