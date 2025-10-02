import React from 'react';
import {action, Story} from '@ladle/react';
import AppNav from "../components/AppNav";
import UserMenu from "../components/user.menu";

export default {
    title: 'Layout / AppNav',
};

type Props = React.ComponentProps<typeof AppNav>;

const AppNavWrapper: Story<Partial<Props>> = (args) => <AppNav {...args as Props} />;

export const LoggedIn = (args: Props) => <AppNavWrapper {...args} />
LoggedIn.args = {
    rightContent: (
        <UserMenu 
            user={{email: "foo", name: "bar"}} 
            onLogout={action('onLogout')} 
        />
    )
};

export const LoggedOut = (args: Props) => <AppNavWrapper {...args} />
LoggedOut.args = {};