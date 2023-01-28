import {Alert, Container, Paper, Toolbar} from "@mui/material";
import AppNav from "../../components/AppNav";
import React from "react";

export default function ErrorPage({error}: {
    error: Error
}) {
    return (
        <>
            <AppNav rightContent={<></>}/>
            <Toolbar/>
            <Container>
                <Paper elevation={3} sx={theme => ({
                    maxWidth: 'sm',
                    margin: theme.spacing(4, 'auto'),
                })}>
                    <Alert severity='error'>{error.message}</Alert>
                </Paper>
            </Container>
        </>
    )
}