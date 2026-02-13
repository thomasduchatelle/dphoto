'use client';

import {Avatar, AvatarGroup, Tooltip} from '@mui/material';
import {UserAvatar} from '@/components/UserAvatar';

export interface SharedByIndicatorProps {
    users: Array<{
        name: string;
        email: string;
        picture?: string;
    }>;
    maxVisible?: number;
}

export const SharedByIndicator = ({users, maxVisible = 3}: SharedByIndicatorProps) => {
    const visibleUsers = users.slice(0, maxVisible);
    const overflowCount = users.length - maxVisible;
    const allUserNames = users.map(u => u.name).join(', ');

    return (
        <AvatarGroup
            max={maxVisible + 1}
            spacing={-8}
            sx={{
                '& .MuiAvatar-root': {
                    border: '1px solid rgba(255, 255, 255, 0.2)',
                },
            }}
            aria-label={`Shared with ${users.length} users: ${allUserNames}`}
        >
            {visibleUsers.map((user, index) => (
                <UserAvatar key={index} name={user.name} picture={user.picture} size="small"/>
            ))}
            {overflowCount > 0 && (
                <Tooltip title={users.slice(maxVisible).map(u => u.name).join(', ')} arrow>
                    <Avatar
                        sx={{
                            bgcolor: 'primary.main',
                            color: '#ffffff',
                            width: 32,
                            height: 32,
                            fontSize: 14,
                            cursor: 'help',
                        }}
                        tabIndex={0}
                    >
                        +{overflowCount}
                    </Avatar>
                </Tooltip>
            )}
        </AvatarGroup>
    );
};
