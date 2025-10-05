'use client';

import {memo, useState} from "react";
import {Avatar, AvatarGroup, Box, Button, Menu, MenuItem} from "@mui/material";
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import {AlbumFilterCriterion, AlbumFilterEntry} from "../../../core/catalog";

export interface OwnerSelectorProps {
    selected: AlbumFilterEntry
    options: AlbumFilterEntry[]
    onAlbumFiltered: (criterion: AlbumFilterCriterion) => void
}

const Avatars = memo(({avatars}: { avatars: string[] }) => {
    return (
        <AvatarGroup max={4} spacing='small' variant="circular" sx={{
            '& .MuiAvatarGroup-avatar': {width: 32, height: 32, fontSize: "0.8em"},
        }}>
            {avatars.length === 0 && (
                <Avatar alt="Empty" title="Empty"/>
            )}
            {avatars.map(avatar => (
                <Avatar key={avatar} alt={avatar} src={avatar}/>
            ))}
        </AvatarGroup>
    )
})

export function OwnerSelector({selected, options = [], onAlbumFiltered}: OwnerSelectorProps) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const handleClickListItem = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleClickOnOption = (criterion: AlbumFilterCriterion) => {
        setAnchorEl(null);
        onAlbumFiltered(criterion)
    }

    const disabled = !options || options.length <= 1;
    return (
        <div>
            <Button
                variant="outlined"
                onClick={handleClickListItem}
                disabled={disabled}
                startIcon={
                    <Avatars avatars={selected.avatars}/>
                }
                endIcon={<ArrowDropDownIcon/>}
            >
                {selected.name}
            </Button>
            <Menu
                id="lock-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    'aria-labelledby': 'lock-button',
                    role: 'listbox',
                }}
            >
                {options.map((option, index) => (
                    <MenuItem key={option.name} divider={option.criterion.selfOwned && option.criterion.owners.length >= 1}
                              onClick={() => handleClickOnOption(option.criterion)}>
                        <Box sx={{mr: 1, width: '90px'}}>
                            <Avatars avatars={option.avatars}/>
                        </Box>
                        {option.name}
                    </MenuItem>
                ))}
            </Menu>
        </div>
    );
}