'use client';

import {
    Avatar,
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Grid,
    IconButton,
    InputAdornment,
    Stack,
    TextField,
    Tooltip,
    useMediaQuery,
    useTheme
} from "@mui/material";
import React, {useRef, useState} from "react";
import "./ShareDialogChipsAnimation.css";
import {ShareError, Sharing, UserDetails} from "../../../core/catalog";
import {Add as AddIcon, Check as CheckIcon, Delete, ErrorOutline as ErrorOutlineIcon, Send as SendIcon, Share as ShareIcon} from "@mui/icons-material";

export default function ShareDialog({
                                        open,
                                        sharedWith,
                                        error,
                                        onClose,
                                        onGrant,
                                        onRevoke,
                                        suggestions = [],
                                    }: {
    open: boolean,
    sharedWith: Sharing[],
    error?: ShareError,
    onClose: () => void,
    onGrant: (email: string) => Promise<void>,
    onRevoke: (email: string) => Promise<void>,
    suggestions?: UserDetails[],
}) {
    const [email, setEmail] = useState(error?.type === "grant" ? error.email : "");
    const [recentlyGranted, setRecentlyGranted] = useState<{ [email: string]: boolean }>({});
    const grantTimeouts = useRef<{ [email: string]: NodeJS.Timeout }>({});

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const handleGrantSuggestion = (email: string) => {
        if (recentlyGranted[email]) {
            // Prevent double click
            return;
        }

        setRecentlyGranted(prev => ({...prev, [email]: true}));

        onGrant(email)
            .then(() => {
                if (grantTimeouts.current[email]) {
                    clearTimeout(grantTimeouts.current[email]);
                }
                grantTimeouts.current[email] = setTimeout(() => {
                    setRecentlyGranted(prev => {
                        const copy = {...prev};
                        delete copy[email];
                        return copy;
                    });
                }, 2000);

                setEmail("");
            })
            .catch(() => {
                setEmail(email);
            });
    };

    const savingHandler = () => {
        if (email) {
            handleGrantSuggestion(email)
        }
    }

    React.useEffect(() => {
        const timeouts = grantTimeouts.current;
        return () => {
            Object.values(timeouts).forEach(timeout => clearTimeout(timeout));
        };
    }, []);

    const topSuggestions = suggestions.slice(0, 5);
    const isRecentlyGranted = (email: string) => recentlyGranted[email];

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
                <Grid container spacing={2} alignItems='center'>
                    <Grid size={{xs: 12}}>
                        <TextField
                            autoFocus
                            fullWidth
                            variant={isMobile ? 'standard' : 'outlined'}
                            margin="dense"
                            size='medium'
                            id="email"
                            placeholder="Email Address"
                            type="email"
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEmail(event.target.value)}
                            onKeyDown={(event: React.KeyboardEvent<HTMLInputElement>) => event.key === 'Enter' && savingHandler()}
                            value={email}
                            error={(error && error.email === email)}
                            helperText={error && error.email === email ? error.message : undefined}
                            autoComplete="off"
                            slotProps={{
                                input: {
                                    startAdornment:
                                        <IconButton sx={{pr: '10px', pl: '0'}} aria-label="share" disabled>
                                            <ShareIcon/>
                                        </IconButton>,
                                    endAdornment:
                                        <InputAdornment position="end" variant="filled">
                                            <Tooltip title="Allow this user to see the pictures and videos of your album">
                                                <IconButton sx={{p: '10px'}}
                                                            aria-label="share"
                                                            onClick={savingHandler}
                                                            color="primary">
                                                    <SendIcon/>
                                                </IconButton>
                                            </Tooltip>
                                        </InputAdornment>
                                }
                            }}
                        />
                        {topSuggestions.length > 0 && (
                            <Stack direction="row" spacing={1} flexWrap="wrap">
                                {topSuggestions.map(user => {
                                    const isError = error?.email === user.email;
                                    return (
                                        <Chip
                                            key={user.email}
                                            size="small"
                                            avatar={<Avatar alt={user.name} src={user.picture}/>}
                                            label={user.name}
                                            variant={isError ? "filled" : "outlined"}
                                            color={isError ? "error" : "default"}
                                            onClick={() => handleGrantSuggestion(user.email)}
                                            deleteIcon={isError ? <ErrorOutlineIcon fontSize="small"/> : <AddIcon fontSize="small"/>}
                                            onDelete={() => handleGrantSuggestion(user.email)}
                                            sx={{
                                                cursor: 'pointer',
                                                maxWidth: 180,
                                                overflow: 'hidden',
                                                whiteSpace: 'nowrap',
                                                textOverflow: 'ellipsis',
                                                transition: 'all 0.2s',
                                                fontWeight: isError ? 600 : undefined,
                                            }}
                                            title={`${user.name} <${user.email}>`}
                                        />
                                    );
                                })}
                            </Stack>
                        )}
                    </Grid>
                </Grid>
                {sharedWith.length > 0 && (
                    <Box sx={theme => ({
                        mt: theme.spacing(3),
                    })}>
                        <Box sx={{mb: 1, fontWeight: 600, color: 'text.secondary', fontSize: 15}}>
                            Users with access
                        </Box>
                        <Stack direction="row" spacing={1} flexWrap="wrap">
                            {sharedWith.map(({user}) => {
                                const recentlyGranted = isRecentlyGranted(user.email);
                                const isError = error?.type === "revoke" && error?.email === user.email;
                                return (
                                    <Chip
                                        key={user.email}
                                        avatar={<Avatar alt={user.name} src={user.picture}/>}
                                        label={user.name}
                                        variant="filled"
                                        color={isError ? "error" : recentlyGranted ? "success" : "primary"}
                                        onDelete={() => onRevoke(user.email).catch(async () => {
                                        })}
                                        deleteIcon={
                                            isError
                                                ? <ErrorOutlineIcon/>
                                                : recentlyGranted
                                                    ? <CheckIcon/>
                                                    : <Delete/>
                                        }
                                        sx={{
                                            maxWidth: 180,
                                            overflow: 'hidden',
                                            whiteSpace: 'nowrap',
                                            textOverflow: 'ellipsis',
                                            fontWeight: isError || recentlyGranted ? 600 : undefined,
                                            transition: 'all 0.2s',
                                        }}
                                        title={`${user.name} <${user.email}>`}
                                    />
                                );
                            })}
                        </Stack>
                        {error?.type === "revoke" && (
                            <Box sx={{
                                color: theme => theme.palette.error.main,
                                fontSize: 13,
                                mt: 0.5,
                            }}>
                                {error.message}
                            </Box>
                        )}
                    </Box>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
