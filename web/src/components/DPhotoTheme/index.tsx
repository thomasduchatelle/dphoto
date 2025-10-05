'use client';

import {createTheme, LinkProps, ThemeProvider} from "@mui/material";
import {forwardRef, ReactNode} from "react";

// Simple link component that doesn't use react-router
const LinkBehavior = forwardRef<HTMLAnchorElement,
    { href: string; [key: string]: any }>((props, ref) => {
    const {href, onClick, ...other} = props;
    
    const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        if (onClick) {
            onClick(e);
        }
    };
    
    return <a data-testid="custom-link" ref={ref} href={href} onClick={handleClick} {...other} />;
});

// https://mycolor.space/?hex=%23005792&sub=1
const theme = createTheme({
    palette: {
        primary: {
            light: '#5288C8',
            main: '#005792',
            dark: '#003369',
            contrastText: '#F5FAFF', // #F5FAFF
        },
        secondary: {
            main: '#00B6AA',
        },
        background: {
            default: '#FFFFFF', // #F0F0F0 #F1F1E6
            paper: '#FFFFFF',
        }
    },
    components: {
        MuiLink: {
            defaultProps: {
                component: LinkBehavior
            } as LinkProps
        },
        MuiButtonBase: {
            defaultProps: {
                LinkComponent: LinkBehavior
            }
        }
    }
});

const DPhotoTheme = ({children}: {
    children: ReactNode
}) => (
    <ThemeProvider theme={theme}>
        {children}
    </ThemeProvider>
);

export default DPhotoTheme
