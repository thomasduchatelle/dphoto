'use client';

import Link from 'next/link';
import {Avatar, Box, IconButton, Paper, Typography} from '@mui/material';
import LogoutIcon from '@mui/icons-material/Logout';

export interface UserInfoProps {
    name: string;
    email: string;
    picture?: string;
    logoutUrl: string;
}

export function UserInfo({ name, email, picture, logoutUrl }: UserInfoProps) {
    return (
        <Paper
            sx={{
                position: 'fixed',
                top: 16,
                right: 16,
                display: 'flex',
                alignItems: 'center',
                gap: 1.5,
                borderRadius: '24px',
                px: 2,
                py: 1,
            }}
            elevation={3}
        >
            {picture ? (
                <Avatar
                    alt={name}
                    src={picture}
                    sx={{width: 32, height: 32}}
                />
            ) : (
                <Avatar
                    sx={{
                        width: 32,
                        height: 32,
                        bgcolor: 'action.selected',
                        color: 'text.primary',
                        fontSize: '0.875rem',
                        fontWeight: 600,
                    }}
                >
                    {name.charAt(0).toUpperCase()}
                </Avatar>
            )}
            <Box sx={{display: 'flex', flexDirection: 'column'}}>
                <Typography
                    variant="body2"
                    sx={{
                        fontWeight: 600,
                        lineHeight: 1.2,
                    }}
                >
                    {name}
                </Typography>
                <Typography
                    variant="caption"
                    sx={{
                        color: 'text.secondary',
                        lineHeight: 1.2,
                    }}
                >
                    {email}
                </Typography>
            </Box>
            <IconButton
                href={logoutUrl}
                LinkComponent={Link}
                size="small"
                title="Logout"
                sx={{
                    ml: 0.5,
                    color: 'text.secondary',
                    '&:hover': {
                        bgcolor: 'action.hover',
                    },
                }}
            >
                <LogoutIcon fontSize="small"/>
            </IconButton>
        </Paper>
    );
}
