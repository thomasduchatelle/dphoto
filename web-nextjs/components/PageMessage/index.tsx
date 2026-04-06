'use client';

import {Box, Typography} from '@mui/material';
import {ReactNode} from 'react';

export type PageMessageVariant = 'info' | 'success' | 'error';

export interface PageMessageProps {
    variant?: PageMessageVariant;
    icon: ReactNode;
    title: string;
    message: string;
    children?: ReactNode;
}

const variantColors: Record<PageMessageVariant, { border: string; icon: string }> = {
    info: {
        border: 'rgba(74, 158, 206, 0.5)',
        icon: 'rgba(74, 158, 206, 0.9)',
    },
    success: {
        border: 'rgba(76, 175, 80, 0.5)',
        icon: 'rgba(76, 175, 80, 0.9)',
    },
    error: {
        border: 'rgba(255, 82, 82, 0.5)',
        icon: 'rgba(255, 82, 82, 0.9)',
    },
};

export const PageMessage = ({variant = 'info', icon, title, message, children}: PageMessageProps) => {
    const colors = variantColors[variant];

    return (
        <Box
            sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '60vh',
                px: 3,
            }}
        >
            <Box
                sx={{
                    maxWidth: 900,
                    width: '100%',
                    position: 'relative',
                    py: 8,
                    px: 5,
                    '&::before': {
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        left: '10%',
                        right: '10%',
                        height: '1px',
                        background: `linear-gradient(90deg, transparent, ${colors.border}, transparent)`,
                    },
                    '&::after': {
                        content: '""',
                        position: 'absolute',
                        bottom: 0,
                        left: '10%',
                        right: '10%',
                        height: '1px',
                        background: `linear-gradient(90deg, transparent, ${colors.border}, transparent)`,
                    },
                    '& .MuiButton-contained': {
                        bgcolor: '#185986',
                        color: '#ffffff',
                        px: 4,
                        py: 1.5,
                        textTransform: 'uppercase',
                        letterSpacing: '0.1em',
                        fontSize: '14px',
                        fontWeight: 400,
                        '&:hover': {
                            bgcolor: '#206ba8',
                            boxShadow: '0 0 24px rgba(24, 89, 134, 0.6)',
                        },
                    },
                    '& .MuiButton-outlined': {
                        borderColor: 'rgba(74, 158, 206, 0.4)',
                        color: 'rgba(255, 255, 255, 0.9)',
                        px: 4,
                        py: 1.5,
                        textTransform: 'uppercase',
                        letterSpacing: '0.1em',
                        fontSize: '14px',
                        fontWeight: 400,
                        '&:hover': {
                            borderColor: '#4a9ece',
                            bgcolor: 'rgba(74, 158, 206, 0.1)',
                        },
                    },
                }}
            >
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        textAlign: 'center',
                    }}
                >
                    <Box
                        sx={{
                            fontSize: 64,
                            color: colors.icon,
                            mb: 3,
                            '& > svg': {
                                fontSize: 64,
                            },
                        }}
                    >
                        {icon}
                    </Box>
                    <Typography
                        variant="h4"
                        component="h2"
                        sx={{
                            fontFamily: 'Georgia, serif',
                            fontWeight: 300,
                            mb: 2,
                            color: '#ffffff',
                        }}
                    >
                        {title}
                    </Typography>
                    <Typography
                        variant="body1"
                        sx={{
                            color: 'rgba(255, 255, 255, 0.75)',
                            fontWeight: 300,
                            mb: children ? 4 : 0,
                            maxWidth: 500,
                        }}
                    >
                        {message}
                    </Typography>
                    {children && (
                        <Box sx={{display: 'flex', gap: 2, flexWrap: 'wrap', justifyContent: 'center'}}>
                            {children}
                        </Box>
                    )}
                </Box>
            </Box>
        </Box>
    );
};
