'use client';

import {Avatar} from '@mui/material';
import Image from 'next/image';

export interface UserAvatarProps {
    name: string;
    picture?: string;
    size?: 'small' | 'medium' | 'large';
}

const SIZE_MAP = {
    small: 32,
    medium: 40,
    large: 64,
};

const getInitials = (name: string): string => {
    return name
        .split(' ')
        .map(part => part[0])
        .join('')
        .toUpperCase()
        .substring(0, 2);
};

export const UserAvatar = ({name, picture, size = 'medium'}: UserAvatarProps) => {
    const sizePixels = SIZE_MAP[size];
    const initials = getInitials(name);

    return (
        <Avatar
            alt={name}
            aria-label={`Profile picture for ${name}`}
            sx={{
                width: sizePixels,
                height: sizePixels,
                bgcolor: 'primary.main',
                color: '#ffffff',
                border: '1px solid rgba(255, 255, 255, 0.2)',
            }}
        >
            {picture ? (
                <Image
                    src={picture}
                    alt={name}
                    width={sizePixels}
                    height={sizePixels}
                />
            ) : (
                initials
            )}
        </Avatar>
    );
};
