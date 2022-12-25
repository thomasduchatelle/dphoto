import {createTheme, LinkProps, ThemeProvider} from "@mui/material";
import {forwardRef, ReactNode} from "react";
import {Link as RouterLink, LinkProps as RouterLinkProps} from "react-router-dom";

const LinkBehavior = forwardRef<any,
  Omit<RouterLinkProps, 'to'> & { href: RouterLinkProps['to'] }>((props, ref) => {
  const {href, ...other} = props;
  // Map href (MUI) -> to (react-router)
  return <RouterLink data-testid="custom-link" ref={ref} to={href} {...other} />;
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
      default: '#FFFFFF', //#F1F1E6
      paper: '#DDF2FF',
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
