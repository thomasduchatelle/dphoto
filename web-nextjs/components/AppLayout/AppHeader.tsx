'use client';

import {AppBar, Box, Toolbar} from '@mui/material';
import Link from '@/components/Link';
import Image from 'next/image';
import {AuthenticatedUser} from '@/libs/security/session-service';
import {UserMenu} from './UserMenu';
import {AppNav} from './AppNav';

export interface AppHeaderProps {
    user: AuthenticatedUser;
    logoutUrl: string;
    isScrolled?: boolean;
}

export const AppHeader = ({user, logoutUrl, isScrolled = false}: AppHeaderProps) => {
    return (
        <AppBar
            position="static"
            elevation={0}
            sx={{
                bgcolor: 'transparent',
                backdropFilter: isScrolled ? 'blur(10px)' : 'none',
                transition: 'all 0.3s ease',
                boxShadow: isScrolled ? '0 2px 16px rgba(0, 0, 0, 0.5)' : 'none',
                borderBottom: isScrolled ? '1px solid rgba(74, 158, 206, 0.2)' : 'none',
            }}
        >
            <Toolbar
                sx={{
                    height: {xs: 56, sm: 64},
                }}
            >
                <Box
                    component={Link}
                    href="/"
                    aria-label="Home"
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        textDecoration: 'none',
                        height: 40,
                        mr: {xs: 2, sm: 4},
                        position: 'relative',
                        '&:hover': {
                            opacity: 0.8,
                        },
                    }}
                >
                    <Image
                        src="/dphoto-fulllogo-large.png"
                        alt="DPhoto"
                        width={160}
                        height={40}
                        priority
                        style={{objectFit: 'contain'}}
                    />
                </Box>
                <AppNav/>
                <Box sx={{flexGrow: 1}}/>
                <UserMenu user={user} logoutUrl={logoutUrl}/>
            </Toolbar>
        </AppBar>
    );
};
