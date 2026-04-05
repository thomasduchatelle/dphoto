import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {Box, Button, Typography} from '@mui/material';
import {PhotoLibrary as PhotoLibraryIcon} from '@mui/icons-material';
import {AppBackground} from '@/components/AppLayout/AppBackground';

/**
 * Exploration: NoAlbum Component Design Options
 *
 * This story presents multiple design variations for the NoAlbum empty state.
 * Compare these side-by-side to choose the best visual direction.
 */
const meta = {
    title: 'Exploration/NoAlbum Options',
    parameters: {
        layout: 'padded',
    },
    tags: ['autodocs'],
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

// OPTION 1: Current Implementation (Centered, Simple)
const CurrentImplementation = () => (
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
            <PhotoLibraryIcon/>
        </Box>
        <Typography variant="h5" component="h2" gutterBottom>
            No albums found
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{mb: 3}}>
            Create your first album to get started.
        </Typography>
        <Button variant="contained" onClick={fn()} sx={{bgcolor: 'primary.main'}}>
            Create Album
        </Button>
    </Box>
);

// OPTION 2: Reference Design (Faded borders top/bottom, dialog-like)
const ReferenceDesign = () => (
    <Box
        sx={{
            maxWidth: 600,
            mx: 'auto',
            background: 'rgba(18, 36, 46, 0.6)',
            position: 'relative',
            py: 8,
            px: 5,
            // Top border with gradient fade
            '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: '10%',
                right: '10%',
                height: '1px',
                background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.5), transparent)',
            },
            // Bottom border with gradient fade
            '&::after': {
                content: '""',
                position: 'absolute',
                bottom: 0,
                left: '10%',
                right: '10%',
                height: '1px',
                background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.5), transparent)',
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
                    color: 'rgba(74, 158, 206, 0.4)',
                    mb: 3,
                    '& > svg': {
                        fontSize: 64,
                    },
                }}
            >
                <PhotoLibraryIcon/>
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
                No Albums Found
            </Typography>
            <Typography
                variant="body1"
                sx={{
                    color: 'rgba(255, 255, 255, 0.75)',
                    fontWeight: 300,
                    mb: 4,
                    maxWidth: 400,
                }}
            >
                Create your first album to get started organizing your photos.
            </Typography>
            <Button
                variant="contained"
                onClick={fn()}
                sx={{
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
                }}
            >
                Create Album
            </Button>
        </Box>
    </Box>
);

// OPTION 3: Minimal (No background, just content on page)
const MinimalDesign = () => (
    <Box
        sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            textAlign: 'center',
            maxWidth: 500,
            mx: 'auto',
            py: 10,
        }}
    >
        <Box
            sx={{
                fontSize: 72,
                color: 'rgba(74, 158, 206, 0.3)',
                mb: 3,
                '& > svg': {
                    fontSize: 72,
                },
            }}
        >
            <PhotoLibraryIcon/>
        </Box>
        <Typography
            variant="h4"
            component="h2"
            sx={{
                fontFamily: 'Georgia, serif',
                fontWeight: 300,
                mb: 1.5,
                color: '#ffffff',
            }}
        >
            No Albums Yet
        </Typography>
        <Typography
            variant="body1"
            sx={{
                color: 'rgba(255, 255, 255, 0.6)',
                fontWeight: 300,
                mb: 4,
            }}
        >
            Create your first album to get started
        </Typography>
        <Button
            variant="outlined"
            onClick={fn()}
            sx={{
                borderColor: 'rgba(74, 158, 206, 0.4)',
                color: 'rgba(255, 255, 255, 0.9)',
                px: 4,
                py: 1.5,
                textTransform: 'uppercase',
                letterSpacing: '0.1em',
                fontSize: '13px',
                fontWeight: 300,
                '&:hover': {
                    borderColor: '#4a9ece',
                    bgcolor: 'rgba(74, 158, 206, 0.1)',
                },
            }}
        >
            Create Album
        </Button>
    </Box>
);

// OPTION 4: Card-style with subtle elevation
const CardStyleDesign = () => (
    <Box
        sx={{
            maxWidth: 480,
            mx: 'auto',
            background: 'rgba(30, 30, 30, 0.6)',
            border: '1px solid rgba(74, 158, 206, 0.2)',
            p: 6,
            textAlign: 'center',
            transition: 'all 0.3s ease',
            '&:hover': {
                border: '1px solid rgba(74, 158, 206, 0.4)',
                boxShadow: '0 8px 32px rgba(24, 89, 134, 0.2)',
            },
        }}
    >
        <Box
            sx={{
                fontSize: 56,
                color: 'rgba(74, 158, 206, 0.5)',
                mb: 2.5,
                '& > svg': {
                    fontSize: 56,
                },
            }}
        >
            <PhotoLibraryIcon/>
        </Box>
        <Typography
            variant="h5"
            component="h2"
            sx={{
                fontFamily: 'Georgia, serif',
                fontWeight: 300,
                mb: 1.5,
                color: '#ffffff',
            }}
        >
            No Albums Found
        </Typography>
        <Typography
            variant="body2"
            sx={{
                color: 'rgba(255, 255, 255, 0.7)',
                fontWeight: 300,
                mb: 3.5,
                lineHeight: 1.7,
            }}
        >
            Create your first album to organize and view your photo collection
        </Typography>
        <Button
            variant="contained"
            onClick={fn()}
            sx={{
                bgcolor: '#185986',
                px: 3.5,
                py: 1.25,
                fontSize: '13px',
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                fontWeight: 400,
                '&:hover': {
                    bgcolor: '#206ba8',
                },
            }}
        >
            Create Album
        </Button>
    </Box>
);

// OPTION 5: Monospace/Technical style (matches album metadata style)
const TechnicalDesign = () => (
    <Box
        sx={{
            maxWidth: 520,
            mx: 'auto',
            py: 8,
            px: 5,
            textAlign: 'center',
            position: 'relative',
        }}
    >
        {/* Subtle side accents */}
        <Box
            sx={{
                position: 'absolute',
                left: 0,
                top: '30%',
                bottom: '30%',
                width: '2px',
                background: 'linear-gradient(to bottom, transparent, rgba(74, 158, 206, 0.6), transparent)',
            }}
        />
        <Box
            sx={{
                position: 'absolute',
                right: 0,
                top: '30%',
                bottom: '30%',
                width: '2px',
                background: 'linear-gradient(to bottom, transparent, rgba(74, 158, 206, 0.6), transparent)',
            }}
        />

        <Box
            sx={{
                fontSize: 48,
                color: 'rgba(74, 158, 206, 0.35)',
                mb: 3,
                '& > svg': {
                    fontSize: 48,
                },
            }}
        >
            <PhotoLibraryIcon/>
        </Box>
        <Typography
            variant="h6"
            component="h2"
            sx={{
                fontFamily: '"Courier New", monospace',
                fontWeight: 400,
                textTransform: 'uppercase',
                letterSpacing: '0.15em',
                mb: 2,
                color: 'rgba(255, 255, 255, 0.9)',
                fontSize: '13px',
            }}
        >
            No Albums Found
        </Typography>
        <Typography
            variant="body2"
            sx={{
                fontFamily: '"Courier New", monospace',
                color: 'rgba(255, 255, 255, 0.6)',
                fontWeight: 300,
                mb: 4,
                fontSize: '13px',
                letterSpacing: '0.05em',
            }}
        >
            CREATE FIRST ALBUM TO BEGIN
        </Typography>
        <Button
            variant="outlined"
            onClick={fn()}
            sx={{
                borderColor: '#185986',
                color: '#4a9ece',
                px: 4,
                py: 1.25,
                fontFamily: '"Courier New", monospace',
                textTransform: 'uppercase',
                letterSpacing: '0.15em',
                fontSize: '12px',
                fontWeight: 400,
                '&:hover': {
                    borderColor: '#4a9ece',
                    bgcolor: 'rgba(24, 89, 134, 0.15)',
                    boxShadow: '0 0 16px rgba(24, 89, 134, 0.3)',
                },
            }}
        >
            + Create Album
        </Button>
    </Box>
);

// OPTION 6: Large hero-style (fills more space, dramatic)
const HeroStyleDesign = () => (
    <Box
        sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            textAlign: 'center',
            minHeight: 500,
            py: 10,
            px: 4,
            background: 'linear-gradient(135deg, rgba(24, 89, 134, 0.08) 0%, rgba(24, 89, 134, 0.02) 100%)',
            position: 'relative',
            '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: '20%',
                right: '20%',
                height: '1px',
                background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.4), transparent)',
            },
        }}
    >
        <Box
            sx={{
                fontSize: 96,
                color: 'rgba(74, 158, 206, 0.25)',
                mb: 4,
                opacity: 0.8,
                '& > svg': {
                    fontSize: 96,
                },
            }}
        >
            <PhotoLibraryIcon/>
        </Box>
        <Typography
            variant="h3"
            component="h2"
            sx={{
                fontFamily: 'Georgia, serif',
                fontWeight: 300,
                letterSpacing: '0.08em',
                mb: 2,
                color: '#ffffff',
                textTransform: 'uppercase',
                fontSize: '32px',
            }}
        >
            Your Albums
        </Typography>
        <Typography
            variant="body1"
            sx={{
                color: 'rgba(255, 255, 255, 0.6)',
                fontWeight: 300,
                mb: 5,
                maxWidth: 500,
                fontSize: '16px',
                lineHeight: 1.8,
            }}
        >
            You don't have any albums yet. Create your first album to start organizing your photo memories.
        </Typography>
        <Button
            variant="contained"
            onClick={fn()}
            sx={{
                bgcolor: '#185986',
                color: '#ffffff',
                px: 5,
                py: 1.75,
                textTransform: 'uppercase',
                letterSpacing: '0.12em',
                fontSize: '14px',
                fontWeight: 400,
                boxShadow: '0 4px 16px rgba(24, 89, 134, 0.3)',
                '&:hover': {
                    bgcolor: '#206ba8',
                    boxShadow: '0 6px 24px rgba(24, 89, 134, 0.5)',
                    transform: 'translateY(-2px)',
                },
                transition: 'all 0.3s ease',
            }}
        >
            Create Your First Album
        </Button>
    </Box>
);

export const Option1_Current: Story = {
    render: () => <CurrentImplementation/>,
    parameters: {
        docs: {
            description: {
                story: '**Current implementation** - Simple centered layout with icon, title, message, and button. Clean and minimal.',
            },
        },
    },
};

export const Option2_ReferenceDesign: Story = {
    render: () => <ReferenceDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Reference design from ux-design-direction-final.html** - Dialog-style with faded gradient borders top and bottom, semi-transparent background, Georgia serif typography.',
            },
        },
    },
};

export const Option3_Minimal: Story = {
    render: () => <MinimalDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Minimal approach** - No background container, larger icon, outlined button. Very light and airy.',
            },
        },
    },
};

export const Option4_Card: Story = {
    render: () => <CardStyleDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Card style** - Subtle background with border, hover effect. More traditional card-based approach.',
            },
        },
    },
};

export const Option5_Technical: Story = {
    render: () => <TechnicalDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Technical/Monospace style** - Uses Courier New to match album metadata styling. Vertical accent lines. Outlined button with glow on hover.',
            },
        },
    },
};

export const Option6_Hero: Story = {
    render: () => <HeroStyleDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Hero style** - Large, dramatic layout that fills vertical space. Big icon, emphasized title, gradient background. Makes empty state feel intentional rather than lacking.',
            },
        },
    },
};

// OPTION 7: Full-page AppBackground with horizontal (top/bottom) faded borders
const FullPageHorizontalBorders = () => (
    <AppBackground>
        <Box
            sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100vh',
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
                    // Top border with gradient fade
                    '&::before': {
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        left: '10%',
                        right: '10%',
                        height: '1px',
                        background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.5), transparent)',
                    },
                    // Bottom border with gradient fade
                    '&::after': {
                        content: '""',
                        position: 'absolute',
                        bottom: 0,
                        left: '10%',
                        right: '10%',
                        height: '1px',
                        background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.5), transparent)',
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
                            color: 'rgba(255, 255, 255, 0.75)',
                            mb: 3,
                            '& > svg': {
                                fontSize: 64,
                            },
                        }}
                    >
                        <PhotoLibraryIcon/>
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
                        No Albums Found
                    </Typography>
                    <Typography
                        variant="body1"
                        sx={{
                            color: 'rgba(255, 255, 255, 0.75)',
                            fontWeight: 300,
                            mb: 4,
                            maxWidth: 400,
                        }}
                    >
                        Create your first album to get started organizing your photos.
                    </Typography>
                    <Button
                        variant="contained"
                        onClick={fn()}
                        sx={{
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
                        }}
                    >
                        Create Album
                    </Button>
                </Box>
            </Box>
        </Box>
    </AppBackground>
);

// OPTION 8: Full-page AppBackground with vertical (left/right) faded borders - Technical style
const FullPageVerticalBorders = () => (
    <AppBackground>
        <Box
            sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                minHeight: '100vh',
                px: 3,
            }}
        >
            <Box
                sx={{
                    maxWidth: 520,
                    width: '100%',
                    py: 8,
                    px: 5,
                    textAlign: 'center',
                    position: 'relative',
                }}
            >
                {/* Left vertical gradient border */}
                <Box
                    sx={{
                        position: 'absolute',
                        left: 0,
                        top: '30%',
                        bottom: '30%',
                        width: '2px',
                        background: 'linear-gradient(to bottom, transparent, rgba(74, 158, 206, 0.6), transparent)',
                    }}
                />
                {/* Right vertical gradient border */}
                <Box
                    sx={{
                        position: 'absolute',
                        right: 0,
                        top: '30%',
                        bottom: '30%',
                        width: '2px',
                        background: 'linear-gradient(to bottom, transparent, rgba(74, 158, 206, 0.6), transparent)',
                    }}
                />

                <Box
                    sx={{
                        fontSize: 48,
                        color: 'rgba(74, 158, 206, 0.35)',
                        mb: 3,
                        '& > svg': {
                            fontSize: 48,
                        },
                    }}
                >
                    <PhotoLibraryIcon/>
                </Box>
                <Typography
                    variant="h6"
                    component="h2"
                    sx={{
                        fontFamily: '"Courier New", monospace',
                        fontWeight: 400,
                        textTransform: 'uppercase',
                        letterSpacing: '0.15em',
                        mb: 2,
                        color: 'rgba(255, 255, 255, 0.9)',
                        fontSize: '13px',
                    }}
                >
                    No Albums Found
                </Typography>
                <Typography
                    variant="body2"
                    sx={{
                        fontFamily: '"Courier New", monospace',
                        color: 'rgba(255, 255, 255, 0.6)',
                        fontWeight: 300,
                        mb: 4,
                        fontSize: '13px',
                        letterSpacing: '0.05em',
                    }}
                >
                    CREATE FIRST ALBUM TO BEGIN
                </Typography>
                <Button
                    variant="outlined"
                    onClick={fn()}
                    sx={{
                        borderColor: '#185986',
                        color: '#4a9ece',
                        px: 4,
                        py: 1.25,
                        fontFamily: '"Courier New", monospace',
                        textTransform: 'uppercase',
                        letterSpacing: '0.15em',
                        fontSize: '12px',
                        fontWeight: 400,
                        '&:hover': {
                            borderColor: '#4a9ece',
                            bgcolor: 'rgba(24, 89, 134, 0.15)',
                            boxShadow: '0 0 16px rgba(24, 89, 134, 0.3)',
                        },
                    }}
                >
                    + Create Album
                </Button>
            </Box>
        </Box>
    </AppBackground>
);

export const Option7_FullPageHorizontal: Story = {
    render: () => <FullPageHorizontalBorders/>,
    parameters: {
        layout: 'fullscreen',
        docs: {
            description: {
                story: '**Full-page with horizontal borders** - Uses AppBackground gradient for entire page. Message centered with transparent background and top/bottom faded gradient borders. Similar to Option 2 but fills the whole viewport.',
            },
        },
    },
};

export const Option8_FullPageVertical: Story = {
    render: () => <FullPageVerticalBorders/>,
    parameters: {
        layout: 'fullscreen',
        docs: {
            description: {
                story: '**Full-page with vertical borders** - Uses AppBackground gradient for entire page. Technical/monospace style with left/right vertical gradient borders instead of horizontal. Message centered with transparent background.',
            },
        },
    },
};
