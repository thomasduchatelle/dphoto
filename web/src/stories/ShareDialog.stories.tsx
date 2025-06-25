import React, {useCallback} from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import ShareDialog from "../pages/authenticated/albums/ShareDialog";
import {ShareError, Sharing, UserDetails} from "../core/catalog";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/ShareDialog',
    component: ShareDialog,
} as ComponentMeta<typeof ShareDialog>;

interface InteractiveContainerState {
    sharedWith: Sharing[]
    suggestions: UserDetails[]
    error?: ShareError
    failures?: string[]
}

const InteractiveContainer = ({sharedWith, suggestions, failures = []}: InteractiveContainerState) => {
    const [open, setOpen] = React.useState(true);
    const [state, setState] = React.useState<InteractiveContainerState>({sharedWith, suggestions, failures});

    const onGrant = useCallback(async (email: string) => {
        if (state.failures?.some((forbiddenEmail) => email === forbiddenEmail)) {
            setState(state => ({
                ...state,
                error: {
                    type: "grant",
                    message: `Cannot grant access to ${email}, you are not allowed to do so`,
                    email,
                },
            }));
            throw new Error(`Cannot grant access to ${email}, you are not allowed to do so`);
        }

        setState(({sharedWith, suggestions, failures}) => {
            const index = suggestions.findIndex(u => u.email === email);
            if (index === -1) {
                sharedWith.push({user: {name: email, email}});
            } else {
                const user = suggestions[index];
                sharedWith.push({user});
                suggestions.splice(index, 1);
            }

            sharedWith = sharedWith.sort((a, b) => a.user.name.localeCompare(b.user.name));
            return {sharedWith, suggestions, failures};
        })
    }, [state, setState]);

    const onRevoke = useCallback(async (email: string) => {
        if (state.failures?.some((forbiddenEmail) => email === forbiddenEmail)) {
            setState(state => ({
                ...state,
                error: {
                    type: "revoke",
                    message: `Cannot revoke access from ${email}, you are not allowed to do so`,
                    email,
                },
            }));
            throw new Error(`Cannot revoke access from ${email}, you are not allowed to do so`);
        }

        setState(({sharedWith, suggestions, failures}) => {
            const index = sharedWith.findIndex(s => s.user.email === email);
            if (index >= 0) {
                const user = sharedWith[index].user;
                sharedWith.splice(index, 1);
                suggestions.push(user);
            }
            return {sharedWith, suggestions, failures};
        })
    }, [state, setState]);

    const onClose = useCallback(() => setOpen(false), [setOpen]);

    return (
        <ShareDialog
            open={open}
            sharedWith={state.sharedWith}
            suggestions={state.suggestions}
            error={state.error}
            onGrant={onGrant}
            onRevoke={onRevoke}
            onClose={onClose}
        />
    );
}

const Template: ComponentStory<typeof ShareDialog> = (args) => (<ShareDialog {...args} onGrant={() => Promise.resolve()} onRevoke={() => Promise.resolve()}/>
);
const InteractiveTemplate: ComponentStory<typeof InteractiveContainer> = (args) => (<InteractiveContainer {...args}/>
);

export const Empty = Template.bind({});
Empty.args = {
    open: true,
    sharedWith: [],
};

Empty.parameters = {
    delay: 300,
};
export const WithShares = Template.bind({});
WithShares.args = {
    open: true,
    sharedWith: [
        {user: {name: "Morgan Stark", email: "morgan@stark.com"}},
        {
            user: {name: "Natasha", email: "blackwidow@avenger.com", picture: "black-widow-profile.jpg"},
        },
    ],
};

WithShares.parameters = {
    delay: 300,
};
export const WithRevokeError = Template.bind({});
WithRevokeError.args = {
    open: true,
    error: {
        type: "revoke",
        message: "A generic issue occurred while doing something",
        email: "blackwidow@avenger.com",
    },
    sharedWith: [
        {
            user: {name: "Natasha", email: "blackwidow@avenger.com", picture: "black-widow-profile.jpg"},
        },
        {user: {name: "Tony Stark", email: "tony@stark.com", picture: "tonystark-profile.jpg"}}
    ],
};
WithRevokeError.parameters = {
    delay: 300,
};

export const WithGrantError = Template.bind({});
WithGrantError.args = {
    open: true,
    error: {
        type: "grant",
        message: "A specific issue occurred while setting the role",
        email: "blackwidow@avenger.com",
    },
    suggestions: [
        {name: "Natasha", email: "blackwidow@avenger.com", picture: "black-widow-profile.jpg"}
    ],
    sharedWith: [
        {
            user: {name: "Tony Stark", email: "tony@stark.com", picture: "tonystark-profile.jpg"},
        },
    ],
};
WithGrantError.parameters = {
    delay: 300,
};

export const WithSuggestions = Template.bind({});
WithSuggestions.args = {
    open: true,
    sharedWith: [],
    suggestions: [
        {name: "Tony Stark", email: "tony@stark.com", picture: "tonystark-profile.jpg"},
        {name: "Steve Rogers", email: "steve@avenger.com", picture: "captain-america.jpg"},
        {name: "Natasha Romanoff", email: "natasha@avenger.com", picture: "black-widow-profile.jpg"},
        {name: "Bruce Banner", email: "bruce@avenger.com", picture: "hulk-profile.webp"},
        {name: "Thor Odinson", email: "thor@asgard.com", picture: "thor.jpg"},
        {name: "Clint Barton", email: "clint@avenger.com", picture: "hawkeye.jpg"},
        {name: "Wanda Maximoff", email: "wanda@avenger.com", picture: "scarlet-witch.jpg"},
    ],
};
WithSuggestions.parameters = {
    delay: 300,
};

export const InteractiveChips = InteractiveTemplate.bind({});
InteractiveChips.args = {
    sharedWith: [
        {
            user: {name: "Natasha", email: "natasha@avenger.com", picture: "black-widow-profile.jpg"},
        },
        {user: {name: "Tony Stark", email: "tony@stark.com", picture: "tonystark-profile.jpg"}},
    ],
    suggestions: [
        {name: "Steve Rogers", email: "steve@avenger.com"},
        {name: "Bruce Banner", email: "bruce@avenger.com", picture: "hulk-profile.webp"},
    ],
    failures: ["steve@avenger.com", "tony@stark.com"]
};
InteractiveChips.parameters = {
    delay: 300,
};
