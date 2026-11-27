'use client';

import {Box} from '@mui/material';
import {usePathname} from 'next/navigation';
import Link from '@/components/Link';

export type AppSection = 'albums' | 'tags';

interface NavEntry {
    section: AppSection
    label: string
    href: string
}

const NAV_ENTRIES: NavEntry[] = [
    {section: 'albums', label: 'Albums', href: '/'},
    {section: 'tags', label: 'Tags', href: '/tags'},
];

export function AppNav() {
    const pathname = usePathname();
    const activeSection: AppSection = pathname.startsWith('/tags') ? 'tags' : 'albums';
    return <AppNavView activeSection={activeSection}/>;
}

export interface AppNavViewProps {
    activeSection: AppSection;
}

export function AppNavView({activeSection}: AppNavViewProps) {
    return (
        <Box
            component="nav"
            aria-label="Main navigation"
            sx={{display: 'flex', alignItems: 'center', gap: {xs: 1.5, sm: 3}}}
        >
            {NAV_ENTRIES.map((entry) => {
                const isActive = entry.section === activeSection;
                return (
                    <Box
                        key={entry.section}
                        component={Link}
                        href={entry.href}
                        prefetch={false}
                        aria-current={isActive ? 'page' : undefined}
                        sx={{
                            position: 'relative',
                            textDecoration: 'none',
                            fontSize: {xs: '0.95rem', sm: '1rem'},
                            fontWeight: 600,
                            color: isActive ? '#ffffff' : 'rgba(255, 255, 255, 0.6)',
                            transition: 'color 0.2s ease',
                            '&:hover': {
                                color: '#ffffff',
                            },
                            '&::after': {
                                content: '""',
                                position: 'absolute',
                                left: 0,
                                right: 0,
                                bottom: -6,
                                height: 2,
                                borderRadius: 1,
                                bgcolor: isActive ? 'primary.main' : 'transparent',
                            },
                        }}
                    >
                        {entry.label}
                    </Box>
                );
            })}
        </Box>
    );
}
