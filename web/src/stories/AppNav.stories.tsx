import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';

import AppNav from "../components/AppNav";
import UserMenu from "../components/user.menu";
import {DefaultMenu} from './UserMenu.stories';

export default {
    title: 'Layout/AppNav',
    component: AppNav,
    subcomponents: {UserMenu},
} as ComponentMeta<typeof AppNav>;

const TemplateWithMenu: ComponentStory<typeof AppNav> = (args) => <AppNav
    {...args}
    rightContent={(<DefaultMenu user={{email: "foo", name: "bar"}} onLogout={() => {
    }} {...DefaultMenu.args} />)}
/>;

export const LoggedIn = TemplateWithMenu.bind({});
LoggedIn.args = {}

const Template: ComponentStory<typeof AppNav> = (args) => <AppNav
    {...args}
/>;
export const LoggedOut = Template.bind({});
LoggedOut.args = {};