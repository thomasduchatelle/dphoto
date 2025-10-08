import React from 'react';
import {action, Story} from '@ladle/react';
import AppNav from "../components/AppNav";
import UserMenu from "../components/user.menu";
import {RouterProvider} from "../components/ClientRouter";

export default {
    title: 'Layout / AppNav',
};

type Props = React.ComponentProps<typeof AppNav>;

const AppNavWrapper: Story<Partial<Props>> = (args) => <RouterProvider><AppNav {...args as Props} /></RouterProvider>;

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