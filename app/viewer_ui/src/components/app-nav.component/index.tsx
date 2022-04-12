import {AppBar, Box, Container, Toolbar, useScrollTrigger} from "@mui/material";
import {cloneElement, ReactElement, ReactNode} from "react";

const appVersion = "1.4.0-delta"

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

const AppNapComponent = ({rightContent}: {
  rightContent: ReactNode
}) => (
  <>
    <ElevationScroll>
      <AppBar>
        <Container maxWidth={false}>
          <Toolbar disableGutters>
            <Box component='a' href='/' sx={{flexGrow: 0}}>
              <img src='/dphoto-fulllogo-reversed-50px.png' alt='DPhoto Logo' style={{height: '50px', marginTop: '5px'}} title={"version: " + appVersion} />
            </Box>
            <Box sx={{flexGrow: 1}}/>
            <Box sx={{flexGrow: 0}}>
              {rightContent}
            </Box>
          </Toolbar>
        </Container>
      </AppBar>
    </ElevationScroll>
    <Toolbar/>
  </>
)

export default AppNapComponent
