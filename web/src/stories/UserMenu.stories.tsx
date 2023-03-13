import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import UserMenu from "../components/user.menu";
import {userEvent, within} from "@storybook/testing-library";

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
DefaultMenu.parameters = {
    storyshots: {disable: true},
    delay: 300,
};
DefaultMenu.play = async ({canvasElement}) => {
    const canvas = within(canvasElement);

    await userEvent.click(canvas.getByRole("button"));
};