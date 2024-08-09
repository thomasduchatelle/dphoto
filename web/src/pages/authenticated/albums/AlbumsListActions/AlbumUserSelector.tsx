import {useState} from "react";
import {Avatar, AvatarGroup, Box, Divider, Fab, Menu, MenuItem} from "@mui/material";
import CollectionsIcon from "@mui/icons-material/Collections";

const options = [
    'Show some love to MUI',
    'Show all notification content',
    'Hide sensitive notification content',
    'Hide all notification content',
];

export default function AlarmUserSelector() {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [selectedIndex, setSelectedIndex] = useState(1);
    const open = Boolean(anchorEl);
    const handleClickListItem = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuItemClick = (
        event: React.MouseEvent<HTMLElement>,
        index: number,
    ) => {
        setSelectedIndex(index);
        setAnchorEl(null);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    return (
        <div>
            <Fab variant="extended" onClick={handleClickListItem}>
                <Box sx={{mr: 1, width: '130px'}}>
                    <AvatarGroup max={4} spacing='small'>
                        <Avatar alt="Black Widow" src="/api/static/black-widow-profile.jpg"/>
                        <Avatar alt="Hulk" src="/api/static/hulk-profile.webp"/>
                        <Avatar alt="Tony Stark" src="/api/static/tonystark-profile.jpg"/>
                        <Avatar alt="Agnes Walker" src="/static/images/avatar/4.jpg"/>
                        <Avatar alt="Trevor Henderson" src="/static/images/avatar/5.jpg"/>
                    </AvatarGroup>
                </Box>
                All Albums
                {/*<CollectionsIcon sx={{mr: 1}}/>*/}
                {/*All Albums*/}
            </Fab>
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
                <MenuItem>
                    <Box sx={{mr: 1, width: '130px'}}>
                        <Avatar alt="Tony Stark" src="/api/static/tonystark-profile.jpg"/>
                        {/*<CollectionsIcon/>*/}
                    </Box>
                    My Albums
                </MenuItem>
                <Divider />
                <MenuItem>
                    <Box sx={{mr: 1, width: '130px'}}>
                        <AvatarGroup max={4} spacing='small'>
                            <Avatar alt="Black Widow" src="/api/static/black-widow-profile.jpg"/>
                            <Avatar alt="Hulk" src="/api/static/hulk-profile.webp"/>
                            <Avatar alt="Tony Stark" src="/api/static/tonystark-profile.jpg"/>
                            <Avatar alt="Agnes Walker" src="/static/images/avatar/4.jpg"/>
                            <Avatar alt="Trevor Henderson" src="/static/images/avatar/5.jpg"/>
                        </AvatarGroup>
                    </Box>
                    All Albums
                </MenuItem>
                <MenuItem>
                    <Box sx={{mr: 1, width: '130px'}}>
                        <Avatar alt="Black Widow" src="/api/static/black-widow-profile.jpg"/>
                    </Box>
                    Black Window
                </MenuItem>
                {/*{options.map((option, index) => (*/}
                {/*    <MenuItem*/}
                {/*        key={option}*/}
                {/*        disabled={index === 0}*/}
                {/*        selected={index === selectedIndex}*/}
                {/*        onClick={(event) => handleMenuItemClick(event, index)}*/}
                {/*    >*/}
                {/*        {option}*/}
                {/*    </MenuItem>*/}
                {/*))}*/}
            </Menu>
        </div>
    );
}