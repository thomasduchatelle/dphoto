import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';

import {LoginInternal} from "../pages/Login";
import {initialLoginPageState, reduce} from "../pages/Login/domain/login-reducer";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Pages/Login',
    component: LoginInternal,
    // More on argTypes: https://storybook.js.org/docs/react/api/argtypes
    argTypes: {
        backgroundColor: {control: 'color'},
    },
} as ComponentMeta<typeof LoginInternal>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof LoginInternal> = (args) => <LoginInternal {...args} />;

const stateWhilePrompting = reduce(initialLoginPageState, {type: "OnUnsuccessfulAutoLoginAttempt"});

export const Initial = Template.bind({});

Initial.args = initialLoginPageState;
export const PromptUser = Template.bind({});
PromptUser.args = stateWhilePrompting;

export const OnTimeout = Template.bind({});
OnTimeout.args = {
    timeout: true,
};

export const OnError = Template.bind({});
OnError.args = reduce(stateWhilePrompting, {type: "error", message: "This is an error"});
