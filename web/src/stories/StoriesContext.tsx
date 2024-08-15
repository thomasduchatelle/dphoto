import {Box, BoxProps, createTheme, LinkProps, ThemeProvider} from "@mui/material";
import {forwardRef} from "react";
import {BrowserRouter, Link as RouterLink, LinkProps as RouterLinkProps} from "react-router-dom";

const LinkBehavior = forwardRef<any,
    Omit<RouterLinkProps, 'to'> & { href: RouterLinkProps['to'] }>((props, ref) => {
    const {href, ...other} = props;
    // Map href (MUI) -> to (react-router)
    return <RouterLink data-testid="custom-link" ref={ref} to={href} {...other} />;
});

// https://mycolor.space/?hex=%23005792&sub=1
const theme = createTheme({
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

export const StoriesContext = (props: BoxProps) => (
    <BrowserRouter>
        <ThemeProvider theme={theme}>
            <Box {...props} />
        </ThemeProvider>
    </BrowserRouter>
);

