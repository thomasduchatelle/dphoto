import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {Box, Button, Typography} from '@mui/material';
import {CloudUpload as CloudUploadIcon} from '@mui/icons-material';

/**
 * Exploration: NoMedia Component Design Options
 *
 * Shows how the chosen NoAlbum design style would apply to NoMedia.
 * These mirror the NoAlbum options but with media-specific content.
 */
const meta = {
    title: 'Exploration/NoMedia Options',
    parameters: {
        layout: 'padded',
    },
    tags: ['autodocs'],
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

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
            '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: '10%',
                right: '10%',
                height: '1px',
                background: 'linear-gradient(90deg, transparent, rgba(74, 158, 206, 0.5), transparent)',
            },
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
                <CloudUploadIcon/>
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
                No Photos Yet
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
                Upload photos to this album to see them displayed here.
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
                Upload Photos
            </Button>
        </Box>
    </Box>
);

// OPTION 5: Monospace/Technical style
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
            <CloudUploadIcon/>
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
            No Photos in Album
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
            Upload Photos to Begin
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
            + Upload Photos
        </Button>
    </Box>
);

// OPTION 6: Large hero-style
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
            <CloudUploadIcon/>
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
            Empty Album
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
            This album doesn't contain any photos yet. Upload your first photos to bring this album to life.
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
            Upload Photos
        </Button>
    </Box>
);

export const Option2_ReferenceDesign: Story = {
    render: () => <ReferenceDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Reference design** applied to NoMedia - Dialog-style with faded borders, upload icon.',
            },
        },
    },
};

export const Option5_Technical: Story = {
    render: () => <TechnicalDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Technical/Monospace style** applied to NoMedia - Matches album metadata typography.',
            },
        },
    },
};

export const Option6_Hero: Story = {
    render: () => <HeroStyleDesign/>,
    parameters: {
        docs: {
            description: {
                story: '**Hero style** applied to NoMedia - Large vertical space, dramatic presentation.',
            },
        },
    },
};
