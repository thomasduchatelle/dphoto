import React from 'react';
import {action, Story} from '@ladle/react';
import UserMenu from "../pages/user.menu";

export default {
    title: 'Layout/UserMenu',
};

type Props = React.ComponentProps<typeof UserMenu>;

const UserMenuWrapper: Story<Partial<Props>> = (args) => <UserMenu {...args as Props} />;

export const DefaultMenu = (args: Props) => <UserMenuWrapper {...args} />
DefaultMenu.args = {
    user: {
        name: "Tony Ironman Stark",
        email: "tomy@stark.com",
        picture: "tonystark-profile.jpg"
    },
    onLogout: action('onLogout'),
};