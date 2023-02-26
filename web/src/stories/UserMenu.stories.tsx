import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import UserMenu from "../components/user.menu";

export default {
    title: 'Components/UserMenu',
    component: UserMenu,
} as ComponentMeta<typeof UserMenu>;

const Template: ComponentStory<typeof UserMenu> = (args) => <UserMenu {...args} />;

export const DefaultMenu = Template.bind({});
DefaultMenu.args = {
    user: {
        name: "Tony Ironman Stark",
        email: "tomy@stark.com",
        picture: "/tonystark-profile.jpg"
    },
    onLogout: () => {
    },
};