'use client';

import {Box, Typography} from '@mui/material';
import {Logout as LogoutIcon} from '@mui/icons-material';
import {UserAvatar} from '@/components/UserAvatar';
import {AuthenticatedUser} from '@/libs/security/session-service';
import Link from '@/components/Link';
import {KeyboardEvent, MouseEvent, useEffect, useRef, useState} from 'react';

export interface UserMenuProps {
    user: AuthenticatedUser;
    logoutUrl: string;
}

export const UserMenu = ({user, logoutUrl}: UserMenuProps) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: globalThis.MouseEvent) => {
            if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
                setIsExpanded(false);
            }
        };

        if (isExpanded) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isExpanded]);

    const handleToggle = (event: MouseEvent<HTMLDivElement>) => {
        event.stopPropagation();
        setIsExpanded(!isExpanded);
    };

    const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            setIsExpanded(!isExpanded);
        } else if (event.key === 'Escape') {
            setIsExpanded(false);
        }
    };

    return (
        <Box
            ref={menuRef}
            onClick={handleToggle}
            onKeyDown={handleKeyDown}
            role="button"
            tabIndex={0}
            aria-expanded={isExpanded}
            aria-label="User profile menu"
            sx={{
                display: 'flex',
                alignItems: 'center',
                gap: 0.75,
                borderRadius: '9999px',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                padding: '0.5rem',
                maxWidth: isExpanded ? '320px' : '56px',
                overflow: 'hidden',
                transition: 'all 0.3s ease',
                cursor: 'pointer',
                '&:hover': {
                    maxWidth: '320px',
                    boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
                },
                '&:hover .user-details': {
                    opacity: 1,
                    width: 'auto',
                },
                '&:hover .logout-button': {
                    opacity: 1,
                    width: 'auto',
                },
                ...(isExpanded && {
                    maxWidth: '320px',
                    boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
                    '& .user-details': {
                        opacity: 1,
                        width: 'auto',
                    },
                    '& .logout-button': {
                        opacity: 1,
                        width: 'auto',
                    },
                }),
            }}
        >
            <UserAvatar name={user.name} picture={user.picture} size="medium"/>
            <Box
                className="user-details"
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    width: 0,
                    opacity: 0,
                    transition: 'none',
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                }}
            >
                <Typography
                    variant="body2"
                    sx={{
                        fontSize: '0.875rem',
                        fontWeight: 600,
                        color: '#ffffff',
                        lineHeight: 1.2,
                    }}
                >
                    {user.name}
                </Typography>
                <Typography
                    variant="caption"
                    sx={{
                        fontSize: '0.75rem',
                        color: 'rgba(255, 255, 255, 0.7)',
                        lineHeight: 1.2,
                    }}
                >
                    {user.email}
                </Typography>
            </Box>
            <Box
                component={Link}
                href={logoutUrl}
                className="logout-button"
                prefetch={false}
                onClick={(e: MouseEvent) => e.stopPropagation()}
                sx={{
                    width: 0,
                    opacity: 0,
                    transition: 'none',
                    ml: 0.5,
                    p: 0.5,
                    color: 'rgba(255, 255, 255, 0.7)',
                    borderRadius: '50%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    textDecoration: 'none',
                    '&:hover': {
                        bgcolor: 'rgba(255, 255, 255, 0.1)',
                        color: '#ffffff',
                    },
                }}
                aria-label="Logout"
                title="Logout"
            >
                <LogoutIcon sx={{fontSize: '1.25rem'}}/>
            </Box>
        </Box>
    );
};
