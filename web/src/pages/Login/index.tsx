import {Alert, Box, LinearProgress, Paper, Typography} from "@mui/material";
import React, {ReactNode} from "react";
import useLoginController from "./domain";
import {LoginPageState} from "./domain/login-hook";
import GoogleLoginButton from "./GoogleLoginButton";
import logo from "../../images/dphoto-fulllogo-medium.png"


const Login = ({onSuccessfulAuthentication}: {
    onSuccessfulAuthentication(): void
}) => {
    const ctrl = useLoginController(onSuccessfulAuthentication)
    const {loginWithIdentityToken, onError} = ctrl

    return (
        <LoginInternal {...ctrl}>
            <GoogleLoginButton onError={onError}
                               onIdentitySuccess={loginWithIdentityToken}
            />
        </LoginInternal>
    )
}
export const LoginInternal = ({
                                  error,
                                  loading,
                                  promptForLogin,
                                  stage,
                                  timeout,
                                  children,
                              }: LoginPageState & {
    children?: ReactNode,
}) => {

    return (
        <Box sx={theme => ({
            height: '100vh',
            backgroundColor: '#F0F0F0',
            paddingTop: {
                sm: theme.spacing(10)
            },
        })}>
            <Box sx={(theme) => ({
                margin: 'auto',
                maxWidth: {
                    sm: 'sm',
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
                        <img src={logo} alt='DPhoto Logo'/>
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

                    {promptForLogin && children}
                </Paper>
            </Box>
        </Box>
    )
}

export default Login
