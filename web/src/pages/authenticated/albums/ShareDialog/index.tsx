import {
    Alert,
    Avatar,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    IconButton,
    ListItemIcon,
    ListItemText,
    Menu,
    MenuItem,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField,
    useMediaQuery,
    useTheme
} from "@mui/material";
import React, {useState} from "react";
import Grid from '@mui/material/Unstable_Grid2';
import {ShareError, Sharing} from "../../../../core/catalog";
import {Delete, MoreHoriz} from "@mui/icons-material";

function OptionButton({onRevoke, name, picture}: {
    onRevoke: () => void
    name?: string
    picture?: string
}) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const handleOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget)
    }
    const handleClose = () => {
        setAnchorEl(null)
    }

    const handleRevoke = () => {
        setAnchorEl(null)
        onRevoke()
    }

    return (
        <>
            <IconButton onClick={handleOpen}><MoreHoriz/></IconButton>
            <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                }}
                keepMounted
                open={Boolean(anchorEl)}
                onClose={handleClose}
            >
                {(name || picture) && (
                    <MenuItem>
                        <ListItemIcon>{picture &&
                            <Avatar alt="?" src={picture} sx={{width: '24px', height: '24px'}}/>}</ListItemIcon>
                        <ListItemText>{name}</ListItemText>
                    </MenuItem>
                )}
                <MenuItem
                    onClick={handleRevoke}
                    sx={theme => ({
                        cursor: 'unset',
                        color: theme.palette.error.main,
                        background: theme.palette.error.contrastText,
                    })}>
                    <ListItemIcon sx={{color: 'inherit'}}><Delete/></ListItemIcon>
                    <ListItemText>
                        Revoke
                    </ListItemText>
                </MenuItem>
            </Menu>
        </>
    )
}

export default function ShareDialog({open, sharedWith, error, onClose, onGrant, onRevoke}: {
    open: boolean,
    sharedWith: Sharing[],
    error?: ShareError,
    onClose: () => void,
    onGrant: (email: string) => Promise<void>,
    onRevoke: (email: string) => Promise<void>,
}) {
    const [email, setEmail] = useState("")

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const savingHandler = () => {
        if (email) {
            onGrant(email).then(() => setEmail(""))
        }
    }

    return (
        <Dialog
            open={open}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth='md'
        >
            <DialogTitle>Sharing album to ...</DialogTitle>
            <DialogContent>
                {error?.type === "general" && (
                    <Alert severity='error' sx={theme => ({mb: theme.spacing(2)})}>{error.message}</Alert>
                )}
                <Grid container spacing={2} alignItems='center'>
                    <Grid sm={8} xs={12}>
                        <TextField
                            autoFocus
                            fullWidth
                            variant={isMobile ? 'standard' : 'outlined'}
                            margin="dense"
                            size='medium'
                            id="email"
                            label="Email Address"
                            type="email"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEmail(event.target.value)}
                            onKeyDown={(event: React.KeyboardEvent<HTMLInputElement>) => event.key === 'Enter' && savingHandler()}
                            value={email}
                            error={error?.type === "adding"}
                            helperText={error?.type === "adding" ? error?.message : undefined}
                        />
                    </Grid>
                    <Grid sm={4} xs={12} sx={{
                        textAlign: 'center'
                    }}>
                        <Button
                            variant="contained"
                            size='large'
                            onClick={savingHandler}
                        >
                            Grant access
                        </Button>
                    </Grid>
                </Grid>
                {sharedWith.length > 0 && (
                    <TableContainer component={Box} sx={theme => ({
                        mt: theme.spacing(3),
                    })}>
                        <Table sx={{
                            minWidth: {
                                xs: 0,
                                md: 650,
                            },
                        }} aria-label="simple table" size='small'>
                            <TableHead>
                                <TableRow>
                                    <TableCell width={20} sx={{
                                        display: {
                                            xs: "none",
                                            sm: "table-cell"
                                        },
                                    }}></TableCell>
                                    <TableCell>Email</TableCell>
                                    <TableCell sx={{
                                        display: {
                                            xs: "none",
                                            sm: "table-cell"
                                        },
                                    }}>Name</TableCell>
                                    <TableCell></TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {sharedWith.map(({user}) => (
                                    <TableRow
                                        key={user.email}
                                        sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                    >
                                        <TableCell width={20} sx={{
                                            pl: 0,
                                            pr: 0,
                                            "& .MuiAvatar-root": {
                                                m: 0,
                                            },
                                            display: {
                                                xs: "none",
                                                sm: "table-cell"
                                            },
                                        }}>
                                            <Avatar alt={user.name} src={user.picture}/>
                                        </TableCell>
                                        <TableCell>
                                            {user.email}
                                        </TableCell>
                                        <TableCell sx={{
                                            display: {
                                                xs: "none",
                                                sm: "table-cell"
                                            },
                                        }}>{user.name}</TableCell>
                                        <TableCell>
                                            <OptionButton onRevoke={() => onRevoke(user.email)}
                                                          name={isMobile ? user.name : undefined}
                                                          picture={isMobile ? user.picture : undefined}/>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
