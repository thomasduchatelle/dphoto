'use client';

import {Alert, Box, Container, Paper, Toolbar, Typography} from "@mui/material";
import AppNav from "../AppNav";
import React from "react";

export const ErrorDisplay = ({error}: {
    error: Error
}) => {
    return (
        <>
            <AppNav rightContent={<></>}/>
            <Toolbar/>
            <Container>
                <Paper elevation={3} sx={theme => ({
                    maxWidth: 'md',
                    margin: theme.spacing(4, 'auto'),
                    padding: theme.spacing(3),
                })}>
                    <Alert severity='error' sx={{marginBottom: 2}}>
                        <Typography variant="h6" gutterBottom>
                            An error occurred
                        </Typography>
                        <Typography variant="body1">
                            Something went wrong. Please refresh the page or try again later.
                        </Typography>
                    </Alert>

                    <Box sx={{marginTop: 2}}>
                        <Typography variant="caption" color="text.secondary">
                            Error details have been logged to the console for debugging.
                        </Typography>
                    </Box>
                </Paper>
            </Container>
        </>
    );
};
