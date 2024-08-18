import {memo, useState} from "react";
import {Avatar, AvatarGroup, Box, Button, Menu, MenuItem} from "@mui/material";
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';

export type Owner = string

export interface AlbumFilterCriterion {
    owners: Owner[] // Empty with selfOwned=false means all albums user has access to
    selfOwned?: boolean // Owned by the current user
}

export interface AlbumFilterEntry {
    criterion: AlbumFilterCriterion
    avatars: string[]
    name: string
}

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

    return (
        <div>
            <Button
                variant="outlined"
                onClick={handleClickListItem}
                disabled={!options || options.length === 0}
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
                    <MenuItem key={option.name} divider={option.criterion.selfOwned || option.criterion.owners.length === 0}
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