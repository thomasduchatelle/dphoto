'use client';

import {Badge} from '@mui/material';
import type {SxProps} from '@mui/material/styles';
import ShareIcon from '@mui/icons-material/Share';

export interface ShareBadgeProps {
    count: number;
    color: string;
    sx?: SxProps;
}

export const ShareBadge = ({count, color, sx}: ShareBadgeProps) => (
    <Badge
        badgeContent={count}
        color="primary"
        sx={{
            '& .MuiBadge-badge': {
                fontSize: 9,
                height: 14,
                minWidth: 14,
                padding: '0 4px',
                backgroundColor: color,
                color: '#ffffff',
                fontWeight: 600,
            },
            ...sx,
        }}
    >
        <ShareIcon sx={{fontSize: 16, color: 'rgba(255, 255, 255, 0.7)'}}/>
    </Badge>
);
