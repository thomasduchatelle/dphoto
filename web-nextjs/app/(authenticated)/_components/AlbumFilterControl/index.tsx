'use client';

import {MouseEvent, useState} from 'react';
import {Avatar, AvatarGroup, Button, Menu, MenuItem} from '@mui/material';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import {AlbumFilterEntry} from '@/domains/catalog/language/catalog-state';

export interface AlbumFilterControlProps {
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterEntry) => void;
}

const OwnerAvatars = ({avatars}: { avatars: string[] }) => (
    <AvatarGroup
        max={4}
        spacing="small"
        sx={{
            '& .MuiAvatar-root': {
                width: 28,
                height: 28,
                fontSize: '0.75rem',
                border: '1px solid rgba(255, 255, 255, 0.2)',
            },
        }}
    >
        {avatars.length === 0 && <Avatar/>}
        {avatars.map(avatar => (
            <Avatar key={avatar} src={avatar}/>
        ))}
    </AvatarGroup>
);

export function AlbumFilterControl({filterOptions, activeFilter, onFilterChange}: AlbumFilterControlProps) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const disabled = filterOptions.length <= 1;

    const handleOpen = (event: MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleSelect = (option: AlbumFilterEntry) => {
        setAnchorEl(null);
        onFilterChange(option);
    };

    return (
        <>
            <Button
                variant="outlined"
                size='large'
                onClick={handleOpen}
                disabled={disabled}
                startIcon={<OwnerAvatars avatars={activeFilter.avatars}/>}
                endIcon={<ArrowDropDownIcon/>}
                aria-label="Album filter by owner"
                aria-haspopup="listbox"
                sx={{
                    minHeight: 44,
                    textTransform: 'none',
                    color: 'text.primary',
                    borderColor: 'rgba(255, 255, 255, 0.2)',
                    '&:hover': {
                        borderColor: 'primary.main',
                        backgroundColor: 'rgba(24, 89, 134, 0.1)',
                    },
                    '&:focus-visible': {
                        outline: '2px solid',
                        outlineColor: 'primary.main',
                        outlineOffset: '2px',
                    },
                }}
            >
                {activeFilter.name}
            </Button>
            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
                slotProps={{
                    list: {
                        role: 'listbox',
                        'aria-label': 'Album filter by owner',
                    }
                }}
            >
                {filterOptions.map(option => (
                    <MenuItem
                        key={option.name}
                        selected={option.name === activeFilter.name}
                        onClick={() => handleSelect(option)}
                        sx={{
                            minHeight: 44,
                            gap: 1.5,
                            '&.Mui-selected': {
                                backgroundColor: 'rgba(24, 89, 134, 0.2)',
                            },
                        }}
                    >
                        <OwnerAvatars avatars={option.avatars}/>
                        {option.name}
                    </MenuItem>
                ))}
            </Menu>
        </>
    );
}
