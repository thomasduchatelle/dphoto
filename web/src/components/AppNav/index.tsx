'use client';

import {AppBar, Box, Container, Toolbar, useScrollTrigger} from "@mui/material";
import {cloneElement, ReactElement, ReactNode} from "react";
import fullLogo from "../../images/dphoto-fulllogo-reversed-50px.png"
import shortLogo from "../../images/dphoto-logo-reversed-50px.png"
import {useClientRouter} from "../ClientRouter";

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
}) => {
    const {navigate} = useClientRouter();

    const handleLogoClick = (e: React.MouseEvent) => {
        e.preventDefault();
        navigate('/');
    };

    return (
        <ElevationScroll>
            <AppBar sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
                <Container maxWidth={false}>
                    <Toolbar disableGutters>
                        <Box component='a' href='/' onClick={handleLogoClick} sx={{flexGrow: 0, display: {xs: 'none', lg: 'flex'}, cursor: 'pointer'}}>
                            <img src={fullLogo} alt='DPhoto Logo'
                                 style={{height: '50px', marginTop: '5px'}}
                            />
                        </Box>
                        <Box component='a' href='/' onClick={handleLogoClick} sx={{flexGrow: 0, display: {lg: 'none'}, cursor: 'pointer'}}>
                            <img src={shortLogo} alt='DPhoto Logo'
                                 style={{height: '50px', marginTop: '5px'}}
                            />
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
    );
}

export default AppNav
