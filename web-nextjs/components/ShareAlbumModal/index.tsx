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
    useTheme,
} from '@mui/material';
import {
    Add as AddIcon,
    Check as CheckIcon,
    Delete as DeleteIcon,
    ErrorOutline as ErrorOutlineIcon,
    Send as SendIcon,
    Share as ShareIcon,
} from '@mui/icons-material';
import {useEffect, useRef, useState} from 'react';
import {SharingDialogFrag} from '@/domains/catalog/sharing';
import {ShareError, Sharing, UserDetails} from '@/domains/catalog/language';
import './ShareDialogChipsAnimation.css';

export interface ShareAlbumModalHandlers {
    onClose: () => void;
    onGrant: (email: string) => Promise<void>;
    onRevoke: (email: string) => Promise<void>;
}

export type ShareAlbumModalProps = SharingDialogFrag & ShareAlbumModalHandlers;

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function isValidEmail(email: string): boolean {
    return EMAIL_REGEX.test(email.trim());
}

export function ShareAlbumModal({open, sharedWith, suggestions, error, onClose, onGrant, onRevoke}: ShareAlbumModalProps) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    return (
        <Dialog
            open={open}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth="md"
        >
            {open && (
                <ShareAlbumModalBody
                    isMobile={isMobile}
                    sharedWith={sharedWith}
                    suggestions={suggestions}
                    error={error}
                    onClose={onClose}
                    onGrant={onGrant}
                    onRevoke={onRevoke}
                />
            )}
        </Dialog>
    );
}

interface ShareAlbumModalBodyProps {
    isMobile: boolean;
    sharedWith: Sharing[];
    suggestions: UserDetails[];
    error?: ShareError;
    onClose: () => void;
    onGrant: (email: string) => Promise<void>;
    onRevoke: (email: string) => Promise<void>;
}

function ShareAlbumModalBody({isMobile, sharedWith, suggestions, error, onClose, onGrant, onRevoke}: ShareAlbumModalBodyProps) {
    const [email, setEmail] = useState(error?.type === 'grant' ? error.email : '');
    const [emailTouched, setEmailTouched] = useState(false);
    const [recentlyGranted, setRecentlyGranted] = useState<{[email: string]: boolean}>({});
    const grantTimeouts = useRef<{[email: string]: ReturnType<typeof setTimeout>}>({});

    useEffect(() => {
        const timeouts = grantTimeouts.current;
        return () => {
            Object.values(timeouts).forEach(timeout => clearTimeout(timeout));
        };
    }, []);

    const trimmedEmail = email.trim();
    const invalidFormat = emailTouched && trimmedEmail.length > 0 && !isValidEmail(trimmedEmail);
    const inputError = error && error.email === email ? error.message : undefined;

    const grant = (target: string) => {
        if (recentlyGranted[target]) {
            return;
        }
        setRecentlyGranted(prev => ({...prev, [target]: true}));
        onGrant(target)
            .then(() => {
                if (grantTimeouts.current[target]) {
                    clearTimeout(grantTimeouts.current[target]);
                }
                grantTimeouts.current[target] = setTimeout(() => {
                    setRecentlyGranted(prev => {
                        const copy = {...prev};
                        delete copy[target];
                        return copy;
                    });
                }, 2000);
                setEmail('');
                setEmailTouched(false);
            })
            .catch(() => {
                setEmail(target);
                setRecentlyGranted(prev => {
                    const copy = {...prev};
                    delete copy[target];
                    return copy;
                });
            });
    };

    const grantTyped = () => {
        if (trimmedEmail && isValidEmail(trimmedEmail)) {
            grant(trimmedEmail);
        }
    };

    const topSuggestions = suggestions.slice(0, 5);

    return (
        <>
            <DialogTitle>Sharing album to ...</DialogTitle>
            <DialogContent>
                <Grid container spacing={2} alignItems="center">
                    <Grid size={{xs: 12}}>
                        <TextField
                            autoFocus
                            fullWidth
                            variant={isMobile ? 'standard' : 'outlined'}
                            margin="dense"
                            size="medium"
                            id="share-album-grantee"
                            placeholder="Email Address"
                            type="text"
                            inputMode="email"
                            onChange={(event) => setEmail(event.target.value)}
                            onBlur={() => setEmailTouched(true)}
                            onKeyDown={(event) => {
                                if (event.key === 'Enter') {
                                    setEmailTouched(true);
                                    grantTyped();
                                }
                            }}
                            value={email}
                            error={Boolean(inputError || invalidFormat)}
                            helperText={inputError ?? (invalidFormat ? 'Enter a valid email address' : undefined)}
                            autoComplete="off"
                            inputProps={{
                                'aria-label': 'Email to grant access',
                                autoComplete: 'off',
                                'data-bwignore': true,
                                'data-1p-ignore': true,
                                'data-lpignore': true,
                                'data-form-type': 'other',
                            }}
                            slotProps={{
                                input: {
                                    startAdornment:
                                        <IconButton sx={{pr: '10px', pl: '0'}} aria-label="share" disabled>
                                            <ShareIcon/>
                                        </IconButton>,
                                    endAdornment:
                                        <InputAdornment position="end" variant="filled">
                                            <Tooltip title="Allow this user to see the pictures and videos of your album">
                                                <span>
                                                    <IconButton
                                                        sx={{p: '10px'}}
                                                        aria-label="Grant access"
                                                        onClick={grantTyped}
                                                        disabled={!isValidEmail(trimmedEmail)}
                                                        color="primary"
                                                    >
                                                        <SendIcon/>
                                                    </IconButton>
                                                </span>
                                            </Tooltip>
                                        </InputAdornment>,
                                },
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
                                            variant={isError ? 'filled' : 'outlined'}
                                            color={isError ? 'error' : 'default'}
                                            onClick={() => grant(user.email)}
                                            deleteIcon={isError ? <ErrorOutlineIcon fontSize="small"/> : <AddIcon fontSize="small"/>}
                                            onDelete={() => grant(user.email)}
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
                    <Box sx={theme => ({mt: theme.spacing(3)})}>
                        <Box sx={{mb: 1, fontWeight: 600, color: 'text.secondary', fontSize: 15}}>
                            Users with access
                        </Box>
                        <Stack direction="row" spacing={1} flexWrap="wrap">
                            {sharedWith.map(({user}) => {
                                const granted = recentlyGranted[user.email];
                                const isError = error?.type === 'revoke' && error?.email === user.email;
                                return (
                                    <Chip
                                        key={user.email}
                                        avatar={<Avatar alt={user.name} src={user.picture}/>}
                                        label={user.name}
                                        variant="filled"
                                        color={isError ? 'error' : granted ? 'success' : 'primary'}
                                        onDelete={() => onRevoke(user.email).catch(() => {
                                        })}
                                        deleteIcon={
                                            isError
                                                ? <ErrorOutlineIcon aria-label={`Revoke access for ${user.name || user.email}`}/>
                                                : granted
                                                    ? <CheckIcon/>
                                                    : <DeleteIcon aria-label={`Revoke access for ${user.name || user.email}`}/>
                                        }
                                        sx={{
                                            maxWidth: 180,
                                            overflow: 'hidden',
                                            whiteSpace: 'nowrap',
                                            textOverflow: 'ellipsis',
                                            fontWeight: isError || granted ? 600 : undefined,
                                            transition: 'all 0.2s',
                                        }}
                                        title={`${user.name} <${user.email}>`}
                                    />
                                );
                            })}
                        </Stack>
                        {error?.type === 'revoke' && (
                            <Box sx={{color: theme => theme.palette.error.main, fontSize: 13, mt: 0.5}}>
                                {error.message}
                            </Box>
                        )}
                    </Box>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </>
    );
}
