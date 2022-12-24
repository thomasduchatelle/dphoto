import {AppBar, Box, Container, Toolbar, useScrollTrigger} from "@mui/material";
import {cloneElement, ReactElement, ReactNode} from "react";

const appVersion = "2.1.0"

const ElevationScroll = ({children}: {
  children: ReactElement;
}) => {
  const trigger = useScrollTrigger({
    disableHysteresis: true,
    threshold: 0,
  });

  return cloneElement(children, {
    elevation: trigger ? 4 : 0,
  });
}

const AppNav = ({rightContent, middleContent}: {
  rightContent: ReactNode,
  middleContent?: ReactNode
}) => (
  <ElevationScroll>
    <AppBar sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
      <Container maxWidth={false}>
        <Toolbar disableGutters>
          <Box component='a' href='/' sx={{flexGrow: 0, display: {xs: 'none', lg: 'flex'}}}>
            <img src='/dphoto-fulllogo-reversed-50px.png' alt='DPhoto Logo' style={{height: '50px', marginTop: '5px'}}
                 title={"version: " + appVersion}/>
          </Box>
          <Box component='a' href='/' sx={{flexGrow: 0, display: {lg: 'none'}}}>
            <img src='/dphoto-logo-reversed-50px.png' alt='DPhoto Logo' style={{height: '50px', marginTop: '5px'}}
                 title={"version: " + appVersion}/>
          </Box>
          <Box sx={{flexGrow: 1}}>
            {middleContent}
          </Box>
          <Box sx={{flexGrow: 0}}>
            {rightContent}
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  </ElevationScroll>
)

export default AppNav
