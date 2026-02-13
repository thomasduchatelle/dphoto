'use client';

import {Box} from '@mui/material';
import {ReactNode, useEffect, useState} from 'react';
import {AppHeader} from '@/components/AppHeader';
import {AppBackground} from './AppBackground';
import {AuthenticatedUser} from '@/libs/security/session-service';

export interface AppLayoutProps {
    children: ReactNode;
    user: AuthenticatedUser;
    logoutUrl: string;
    basePath?: string;
}

export default function AppLayout({children, user, logoutUrl, basePath}: AppLayoutProps) {
    const [isScrolled, setIsScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            setIsScrolled(window.scrollY > 10);
        };

        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    return (
        <AppBackground>
            <Box
                component="header"
                sx={{
                    position: 'fixed',
                    top: 0,
                    left: 0,
                    right: 0,
                    zIndex: 1100,
                }}
            >
                <AppHeader user={user} logoutUrl={logoutUrl} isScrolled={isScrolled} basePath={basePath}/>
            </Box>
            <Box
                component="main"
                sx={{
                    marginTop: {xs: '56px', sm: '64px'},
                    padding: {xs: 2, sm: 3, md: 4},
                    flexGrow: 1,
                }}
            >
                {children}
            </Box>
        </AppBackground>
    );
};
