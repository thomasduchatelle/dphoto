import type {Meta, StoryObj} from '@storybook/react';
import AppLayout from './index';
import {Box, Typography} from '@mui/material';
import {AppHeader} from '@/components/AppHeader';

const meta = {
    title: 'Layout/AppLayout',
    component: AppLayout,
    parameters: {
        layout: 'fullscreen',
    },
    args: {
        user: {
            name: 'Tony Stark',
            email: 'tony@stark-industries.com',
            picture: '/tonystark-profile.jpg',
        },
        logoutUrl: '/auth/logout',
        basePath: "",
        children: (
            <Box>
                <Typography variant="h1" sx={{mb: 2}}>
                    Page Title
                </Typography>
                <Typography variant="body1" sx={{mb: 6}}>
                    Scroll down to see the header blur effect activate.
                </Typography>
                {[1, 2, 3, 4].map((i) => (
                    <Box key={i}>
                        <Typography variant="h2" sx={{mb: 3}}>
                            Section {i}
                        </Typography>
                        <Box
                            sx={{
                                mb: 6,
                                height: 400,
                                bgcolor: 'rgba(255,255,255,0.05)',
                                borderRadius: 1,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                color: 'rgba(255,255,255,0.3)',
                            }}
                        >
                            Image Placeholder
                        </Box>
                    </Box>
                ))}
            </Box>
        ),
    },
} satisfies Meta<typeof AppLayout>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Scrolled: Story = {
    render: (args) => (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                flexDirection: 'column',
                background: 'linear-gradient(135deg, #001929 0%, #0a1520 25%, #12242e 50%, #0f1d28 75%, #001929 100%)',
                backgroundAttachment: 'fixed',
            }}
        >
            <Box component="header" sx={{position: 'fixed', top: 0, left: 0, right: 0, zIndex: 1100}}>
                <AppHeader user={args.user} logoutUrl={args.logoutUrl} isScrolled={true} basePath=""/>
            </Box>
            <Box component="main" sx={{marginTop: 0, padding: {xs: 2, sm: 3, md: 4}, flexGrow: 1}}>
                <Box
                    sx={{
                        height: 200,
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        borderRadius: 1,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: 'white',
                        mb: 4,
                    }}
                >
                    Notice the blur effect on the header above
                </Box>
                {args.children}
            </Box>
        </Box>
    ),
};

export const Mobile: Story = {
    globals: {
        viewport: {
            value: 'mobile2',
            isRotated: false,
        },
    },
};
