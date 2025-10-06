import React from 'react';
import {Story} from '@ladle/react';

import {Button} from "@mui/material";
import {LoginInternal} from "../pages-old/Login";
import {initialLoginPageState, reduce} from 'src/pages-old/Login/domain/login-reducer';

export default {
    title: 'Layout/Login',
};

type LoginProps = React.ComponentProps<typeof LoginInternal>;

const LoginWrapper: Story<Partial<LoginProps>> = (props) => (
    <LoginInternal {...props as LoginProps}>
        <Button sx={{
            width: '100%',
            maxWidth: '400px',
            marginTop: '24px',
        }} variant='outlined' color='primary'>SSO Placeholder</Button>
    </LoginInternal>
);

const stateWhilePrompting = reduce(initialLoginPageState, {type: "OnUnsuccessfulAutoLoginAttempt"});

export const Loading = (args: LoginProps) => <LoginWrapper {...args} />
Loading.args = initialLoginPageState;
Loading.meta = {skipSnapshot: true}

export const PromptUser = (args: LoginProps) => <LoginWrapper {...args} />
PromptUser.args = stateWhilePrompting;

export const OnTimeout = (args: LoginProps) => <LoginWrapper {...args} />
OnTimeout.args = reduce(initialLoginPageState, {type: "OnExpiredSession"});

export const OnError = (args: LoginProps) => <LoginWrapper {...args} />
OnError.args = reduce(stateWhilePrompting, {type: "error", message: "This is an error"});
