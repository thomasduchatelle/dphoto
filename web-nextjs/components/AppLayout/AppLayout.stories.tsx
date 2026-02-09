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
    tags: ['autodocs'],
    args: {
        user: {
            name: 'Tony Stark',
            email: 'tony@stark-industries.com',
            picture: '/tonystark-profile.jpg',
        },
        logoutUrl: '/auth/logout',
        basePath: "",
    },
    globals: {
        viewport: {},
    },
} satisfies Meta<typeof AppLayout>;

export default meta;
type Story = StoryObj<typeof meta>;

const NeutralContent = () => (
    <Box>
        <Typography variant="h1" sx={{mb: 2}}>
            Page Title Example
        </Typography>
        <Typography variant="body1" sx={{mb: 6}}>
            This is an example of body text. The layout includes a responsive header with logo,
            user menu, and smooth scrolling behavior. Scroll down to see the header become blurred
            with a shadow when content passes underneath.
        </Typography>

        <Typography variant="h2" sx={{mb: 3}}>
            Section One
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
                fontSize: '0.875rem',
            }}
        >
            Image Placeholder
        </Box>

        <Typography variant="h2" sx={{mb: 3}}>
            Section Two
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
                fontSize: '0.875rem',
            }}
        >
            Image Placeholder
        </Box>

        <Typography variant="h2" sx={{mb: 3}}>
            Section Three
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
                fontSize: '0.875rem',
            }}
        >
            Image Placeholder
        </Box>

        <Typography variant="h2" sx={{mb: 3}}>
            Section Four
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
                fontSize: '0.875rem',
            }}
        >
            Image Placeholder
        </Box>
    </Box>
);

/**
 * Default state showing the layout with scrollable content.
 * The header is transparent at the top. Scroll down to see the header blur effect activate automatically.
 */
export const Default: Story = {
    args: {
        children: <NeutralContent/>,
    },
};

/**
 * Shows the header in scrolled state with blur effect, shadow, and bottom border.
 * This demonstrates how the header looks when content passes underneath it.
 *
 * Note: This is a static demonstration. In the real application (Default story),
 * this effect activates automatically when scrolling.
 */
export const Scrolled: Story = {
    render: (args) => {
        // Create a static demo showing the scrolled state
        return (
            <Box
                sx={{
                    minHeight: '100vh',
                    display: 'flex',
                    flexDirection: 'column',
                    background: 'linear-gradient(135deg, #001929 0%, #0a1520 25%, #12242e 50%, #0f1d28 75%, #001929 100%)',
                    backgroundAttachment: 'fixed',
                }}
            >
                {/* Fixed header in scrolled state */}
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
                    <AppHeader user={args.user} logoutUrl={args.logoutUrl} isScrolled={true}/>
                </Box>

                {/* Content with visual elements that show behind the blurred header */}
                <Box
                    component="main"
                    sx={{
                        marginTop: 0, // No margin to show content under header
                        padding: {xs: 2, sm: 3, md: 4},
                        flexGrow: 1,
                    }}
                >
                    {/* Colorful content to show the blur effect */}
                    <Box
                        sx={{
                            height: 200,
                            background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                            borderRadius: 1,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            color: 'white',
                            fontSize: '1rem',
                            fontWeight: 300,
                            mb: 4,
                        }}
                    >
                        Notice the blur effect on the header above
                    </Box>

                    <Typography variant="body1" sx={{mb: 6, textAlign: 'center', opacity: 0.7}}>
                        The header is transparent but blurs the content underneath.
                        It also has a shadow and bottom border for separation.
                    </Typography>

                    <NeutralContent/>
                </Box>
            </Box>
        );
    },
};

export const DefaultMobile: Story = {
    globals: {
        viewport: {
            value: 'mobile2',
            isRotated: false,
        },
    },
    args: {
        children: <NeutralContent/>,
    },
};

/**
 * Shows the user menu in its hover/expanded state with user details and logout button visible.
 * This demonstrates the hover-expandable user menu component.
 *
 * Note: This is a static demonstration. In the real application (Default story),
 * this effect activates automatically when hovering over the user avatar.
 */
export const UserMenuHoverState: Story = {
    render: (args) => {
        return (
            <Box
                sx={{
                    minHeight: '100vh',
                    display: 'flex',
                    flexDirection: 'column',
                    background: 'linear-gradient(135deg, #001929 0%, #0a1520 25%, #12242e 50%, #0f1d28 75%, #001929 100%)',
                    backgroundAttachment: 'fixed',
                }}
            >
                {/* Fixed header with forced hover state */}
                <Box
                    component="header"
                    sx={{
                        position: 'fixed',
                        top: 0,
                        left: 0,
                        right: 0,
                        zIndex: 1100,
                        // Force hover state on the user menu
                        '& [aria-label="User profile"]': {
                            maxWidth: '320px !important',
                            border: '1px solid rgba(74, 158, 206, 0.5) !important',
                            bgcolor: 'rgba(24, 89, 134, 0.15) !important',
                            boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05) !important',
                            '& .user-details': {
                                opacity: '1 !important',
                                width: 'auto !important',
                            },
                            '& .logout-button': {
                                opacity: '1 !important',
                                width: 'auto !important',
                            },
                        },
                    }}
                >
                    <AppHeader user={args.user} logoutUrl={args.logoutUrl} isScrolled={false}/>
                </Box>

                {/* Content */}
                <Box
                    component="main"
                    sx={{
                        marginTop: {xs: '56px', sm: '64px'},
                        padding: {xs: 2, sm: 3, md: 4},
                        flexGrow: 1,
                    }}
                >
                    <Typography variant="h2" sx={{mb: 3, textAlign: 'center'}}>
                        User Menu Hover State
                    </Typography>
                    <Typography variant="body1" sx={{mb: 6, textAlign: 'center', opacity: 0.7}}>
                        The user menu in the header is shown in its expanded hover state,
                        displaying the avatar, name, email, and logout button.
                    </Typography>

                    <NeutralContent/>
                </Box>
            </Box>
        );
    },
};
