import {createTheme} from '@mui/material/styles';

export const theme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#185986',
        },
        background: {
            default: '#121212',
            paper: '#1e1e1e',
        },
        text: {
            primary: '#ffffff',
            secondary: 'rgba(255, 255, 255, 0.7)',
        },
    },
    breakpoints: {
        values: {
            xs: 0,
            sm: 600,
            md: 960,
            lg: 1280,
            xl: 1920,
        },
    },
    typography: {
        h1: {
            fontSize: '2rem',
            fontWeight: 300,
            color: '#4a9ece',
            letterSpacing: '0.02em',
        },
        h2: {
            fontSize: '13px',
            fontWeight: 400,
            letterSpacing: '0.15em',
            textTransform: 'uppercase',
            color: 'rgba(255, 255, 255, 0.8)',
            position: 'relative',
            paddingBottom: '12px',
            '&::after': {
                content: '""',
                position: 'absolute',
                bottom: 0,
                left: 0,
                width: '60px',
                height: '1px',
                background: 'rgba(74, 158, 206, 0.6)',
            },
        },
        body1: {
            fontSize: '1rem',
            color: 'rgba(255, 255, 255, 0.7)',
            lineHeight: 1.6,
        },
    },
});
