import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import {Box, IconButton} from "@mui/material";
import React from "react";
import {Link} from "react-router-dom";

export function FullHeightLink({mediaLink, side}: {
  mediaLink: string | undefined
  side: 'left' | 'right'
}) {
  if (!mediaLink) {
    return null
  }
  return (
    <Box sx={theme => ({
      position: 'absolute',
      top: 0,
      ...(side === 'left' && {left: 0}),
      ...(side === 'right' && {right: 0}),
      height: '100vh',
      width: '25%',
      color: theme.palette.background.paper,
      '& > a': {
        display: 'flex',
        width: '100%',
        height: '100%',
        padding: theme.spacing(4),
        alignItems: 'center',
        justifyContent: side,
        color: 'inherit',
      },
      '& button': {
        color: 'inherit',
        display: 'none',
      },
      [theme.breakpoints.up("lg")]: {
        '& button.Mui-focusVisible, & a:hover button': {
          backgroundColor: 'rgba(122, 122, 122, 0.3)',
          display: 'flex',
        },
      },
    })}>
      <Link to={mediaLink}>
        <Box sx={{display: 'flex'}}>
          <IconButton size='large'>
            {side === "left" ? (
              <ChevronLeftIcon fontSize='large'/>
            ) : (
              <ChevronRightIcon fontSize='large'/>
            )}
          </IconButton>
        </Box>
      </Link>
    </Box>
  );
}