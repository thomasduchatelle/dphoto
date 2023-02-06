import {Alert, Box, LinearProgress, Paper, Typography} from "@mui/material";
import React from "react";
import GoogleLoginIntegration from "./GoogleLoginButton";
import useLoginController from "./domain";


const Login = ({onSuccessfulAuthentication}: {
    onSuccessfulAuthentication(): void
}) => {
    const {error, loading, loginWithIdentityToken, onError, onWaitingForUserInput, stage, timeout} = useLoginController(onSuccessfulAuthentication)

    return (
        <Box sx={{
            height: '100vh',
        }}>
            <Box sx={(theme) => ({
                margin: 'auto',
                maxWidth: {
                    sm: 'sm',
                },
                marginTop: {
                    sm: theme.spacing(10),
                },
            })}>
                <Paper elevation={3} sx={(theme) => ({
                    textAlign: 'center',
                    paddingBottom: theme.spacing(3),
                    '& a, & > div, & h4': {
                        marginTop: theme.spacing(3),
                    },
                    height: {
                        xs: '100vh',
                        sm: null,
                    },
                })}>
                    <Box sx={{
                        height: '4px',
                        marginTop: '0px !important',
                    }}>
                        {loading && <LinearProgress sx={{
                            borderRadius: {
                                sm: '4px 4px 0px 0px'
                            },
                        }}/>}
                    </Box>

                    <Box component='a' href='/' sx={theme => ({
                        display: 'block',
                        paddingTop: theme.spacing(3).sub(),
                    })}>
                        <img src='/dphoto-fulllogo-medium.png' alt='DPhoto Logo'/>
                    </Box>

                    <Box>
                        <Typography variant='body1'>This is an invitation only application. Sign in with your Google
                            account.</Typography>
                    </Box>

                    {stage && (
                        <Box>
                            <Typography variant='caption'>{stage}</Typography>
                        </Box>
                    )}
                    {error && (
                        <Alert severity='error'>{error}</Alert>
                    )}
                    {timeout && (
                        <Alert severity={"warning"}>Your session has timed out, thank you to reconnect</Alert>
                    )}

                    <GoogleLoginIntegration onError={onError}
                                            onIdentitySuccess={loginWithIdentityToken}
                                            onWaitingUserInput={onWaitingForUserInput}
                    />
                </Paper>
            </Box>
        </Box>
    )
}

export default Login
