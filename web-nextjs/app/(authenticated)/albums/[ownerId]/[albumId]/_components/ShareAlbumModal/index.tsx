'use client';

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
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    Stack,
    TextField,
    Typography,
    useMediaQuery,
    useTheme,
} from '@mui/material';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import ShareIcon from '@mui/icons-material/Share';
import {useState} from 'react';
import {SharingDialogFrag} from '@/domains/catalog/sharing';
import {ShareError, Sharing} from '@/domains/catalog/language';

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

export function ShareAlbumModal({open, sharedWith, error, onClose, onGrant, onRevoke}: ShareAlbumModalProps) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    return (
        <Dialog
            open={open}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth="sm"
            aria-labelledby="share-album-modal-title"
        >
            <DialogTitle id="share-album-modal-title">
                <Stack direction="row" alignItems="center" gap={1}>
                    <ShareIcon fontSize="small"/>
                    <span>Share album</span>
                </Stack>
            </DialogTitle>
            {open && (
                <ShareAlbumModalBody
                    sharedWith={sharedWith}
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
    sharedWith: Sharing[];
    error?: ShareError;
    onClose: () => void;
    onGrant: (email: string) => Promise<void>;
    onRevoke: (email: string) => Promise<void>;
}

function ShareAlbumModalBody({sharedWith, error, onClose, onGrant, onRevoke}: ShareAlbumModalBodyProps) {
    const initialEmail = error?.type === 'grant' ? error.email : '';
    const [email, setEmail] = useState(initialEmail);
    const [emailTouched, setEmailTouched] = useState(initialEmail.length > 0);
    const [granting, setGranting] = useState(false);
    const [pendingRevoke, setPendingRevoke] = useState<string | undefined>(undefined);
    const [revoking, setRevoking] = useState(false);

    const trimmedEmail = email.trim();
    const invalidFormat = emailTouched && trimmedEmail.length > 0 && !isValidEmail(trimmedEmail);
    const grantApiError = error?.type === 'grant' && error.email === trimmedEmail ? error.message : undefined;
    const grantDisabled = !isValidEmail(trimmedEmail) || granting;

    const handleGrant = () => {
        if (grantDisabled) return;
        setGranting(true);
        onGrant(trimmedEmail)
            .then(() => {
                setEmail('');
                setEmailTouched(false);
            })
            .catch(() => {
            })
            .finally(() => {
                setGranting(false);
            });
    };

    const handleRevokeConfirm = () => {
        if (!pendingRevoke) return;
        const target = pendingRevoke;
        setRevoking(true);
        onRevoke(target)
            .catch(() => {
            })
            .finally(() => {
                setRevoking(false);
                setPendingRevoke(undefined);
            });
    };

    const pendingRevokeSharing = pendingRevoke ? sharedWith.find(s => s.user.email === pendingRevoke) : undefined;
    const revokeApiError = error?.type === 'revoke' ? error : undefined;

    return (
        <>
            <DialogContent dividers>
                {!pendingRevokeSharing ? (
                    <Stack gap={3}>
                        <Box>
                            <Typography variant="body2" sx={{mb: 1, color: 'text.secondary'}}>
                                Grant access by email
                            </Typography>
                            <Stack direction={{xs: 'column', sm: 'row'}} gap={1} alignItems={{xs: 'stretch', sm: 'flex-start'}}>
                                <TextField
                                    fullWidth
                                    autoFocus
                                    size="small"
                                    type="email"
                                    placeholder="user@example.com"
                                    value={email}
                                    onChange={(event) => setEmail(event.target.value)}
                                    onBlur={() => setEmailTouched(true)}
                                    onKeyDown={(event) => {
                                        if (event.key === 'Enter') {
                                            setEmailTouched(true);
                                            handleGrant();
                                        }
                                    }}
                                    error={Boolean(invalidFormat || grantApiError)}
                                    helperText={grantApiError ?? (invalidFormat ? 'Enter a valid email address' : ' ')}
                                    inputProps={{'aria-label': 'Email to grant access'}}
                                    disabled={granting}
                                />
                                <Button
                                    variant="contained"
                                    onClick={handleGrant}
                                    disabled={grantDisabled}
                                    sx={{alignSelf: {xs: 'flex-end', sm: 'flex-start'}, minWidth: 96}}
                                >
                                    {granting ? 'Granting…' : 'Grant'}
                                </Button>
                            </Stack>
                            {grantApiError && (
                                <Alert
                                    severity="error"
                                    sx={{mt: 1}}
                                    action={
                                        <Button color="inherit" size="small" onClick={handleGrant} disabled={grantDisabled}>
                                            Retry
                                        </Button>
                                    }
                                >
                                    {grantApiError}
                                </Alert>
                            )}
                        </Box>
                        <Box>
                            <Typography variant="body2" sx={{mb: 1, color: 'text.secondary'}}>
                                Users with access
                            </Typography>
                            {sharedWith.length === 0 ? (
                                <Typography variant="body2" sx={{color: 'text.secondary', fontStyle: 'italic'}}>
                                    No one has access yet.
                                </Typography>
                            ) : (
                                <List disablePadding>
                                    {sharedWith.map(({user}) => {
                                        const isUserError = revokeApiError?.email === user.email;
                                        return (
                                            <ListItem
                                                key={user.email}
                                                disableGutters
                                                secondaryAction={
                                                    <IconButton
                                                        edge="end"
                                                        aria-label={`Revoke access for ${user.name || user.email}`}
                                                        onClick={() => setPendingRevoke(user.email)}
                                                    >
                                                        <DeleteOutlineIcon/>
                                                    </IconButton>
                                                }
                                            >
                                                <ListItemAvatar>
                                                    <Avatar alt={user.name} src={user.picture}/>
                                                </ListItemAvatar>
                                                <ListItemText
                                                    primary={user.name || user.email}
                                                    secondary={user.name ? user.email : undefined}
                                                    primaryTypographyProps={{color: isUserError ? 'error' : undefined}}
                                                />
                                            </ListItem>
                                        );
                                    })}
                                </List>
                            )}
                            {revokeApiError && (
                                <Alert
                                    severity="error"
                                    sx={{mt: 1}}
                                    action={
                                        <Button
                                            color="inherit"
                                            size="small"
                                            onClick={() => {
                                                setPendingRevoke(revokeApiError.email);
                                            }}
                                        >
                                            Retry
                                        </Button>
                                    }
                                >
                                    {revokeApiError.message}
                                </Alert>
                            )}
                        </Box>
                    </Stack>
                ) : (
                    <Stack gap={2}>
                        <Typography variant="h6">Revoke access?</Typography>
                        <Alert severity="warning">
                            <Typography variant="body2">
                                <strong>{pendingRevokeSharing.user.name || pendingRevokeSharing.user.email}</strong> will no longer be able to view this album.
                            </Typography>
                        </Alert>
                        <Stack direction="row" alignItems="center" gap={1.5}>
                            <Avatar alt={pendingRevokeSharing.user.name} src={pendingRevokeSharing.user.picture}/>
                            <Box sx={{minWidth: 0}}>
                                <Typography noWrap>{pendingRevokeSharing.user.name || pendingRevokeSharing.user.email}</Typography>
                                {pendingRevokeSharing.user.name && (
                                    <Typography variant="body2" color="text.secondary" noWrap>
                                        {pendingRevokeSharing.user.email}
                                    </Typography>
                                )}
                            </Box>
                        </Stack>
                    </Stack>
                )}
            </DialogContent>
            <DialogActions>
                {pendingRevokeSharing ? (
                    <>
                        <Button
                            onClick={() => setPendingRevoke(undefined)}
                            disabled={revoking}
                            autoFocus
                        >
                            Cancel
                        </Button>
                        <Button
                            variant="contained"
                            color="error"
                            onClick={handleRevokeConfirm}
                            disabled={revoking}
                        >
                            {revoking ? 'Revoking…' : 'Revoke'}
                        </Button>
                    </>
                ) : (
                    <Button onClick={onClose}>Close</Button>
                )}
            </DialogActions>
        </>
    );
}
