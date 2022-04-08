import {createTheme, ThemeProvider} from "@mui/material";
import {ReactNode} from "react";

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
      default: '#DDF2FF', //#F1F1E6
      paper: '#FFFFFF',
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
