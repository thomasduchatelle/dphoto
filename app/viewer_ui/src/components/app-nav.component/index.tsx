import {AppBar, Box, Container, Toolbar, useScrollTrigger} from "@mui/material";
import {cloneElement, ReactElement, ReactNode} from "react";

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

export default ({rightContent}: {
  rightContent: ReactNode
}) => (
  <>
    <ElevationScroll>
      <AppBar>
        <Container maxWidth={false}>
          <Toolbar disableGutters>
            <Box component='a' href='/' sx={{flexGrow: 0}}>
              <img src='/dphoto-fulllogo-50px.png' alt='DPhoto' style={{height: '50px', marginTop: '5px'}}/>
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
// <Tooltip title="Open settings">
//   <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
//     <Avatar alt="Remy Sharp" src="/static/images/avatar/2.jpg" />
//   </IconButton>
// </Tooltip>