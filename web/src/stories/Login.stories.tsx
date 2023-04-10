import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';

import {LoginInternal} from "../pages/Login";
import {initialLoginPageState, reduce} from "../pages/Login/domain/login-reducer";
import {Button} from "@mui/material";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Layout/Login',
    component: LoginInternal,
} as ComponentMeta<typeof LoginInternal>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof LoginInternal> = (args) => (<LoginInternal {...args} >
    <Button sx={{
        width: '100%',
        maxWidth: '400px',
        marginTop: '24px',
    }} variant='outlined' color='primary'>SSO Placeholder</Button>
</LoginInternal>);

const stateWhilePrompting = reduce(initialLoginPageState, {type: "OnUnsuccessfulAutoLoginAttempt"});

export const Loading = Template.bind({});

Loading.args = initialLoginPageState;
Loading.parameters = {
    storyshots: {disable: true},
};

export const PromptUser = Template.bind({});
PromptUser.args = stateWhilePrompting;

export const OnTimeout = Template.bind({});
OnTimeout.args = reduce(initialLoginPageState, {type: "OnExpiredSession"});

export const OnError = Template.bind({});
OnError.args = reduce(stateWhilePrompting, {type: "error", message: "This is an error"});
