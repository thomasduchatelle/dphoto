'use client';

import {Alert, Box, Button, Collapse, Typography} from '@mui/material';
import {ErrorOutline as ErrorOutlineIcon} from '@mui/icons-material';
import {useState} from 'react';

export interface ErrorDisplayProps {
    error: {
        message: string;
        code?: string;
        details?: string;
    };
    onRetry?: () => void;
    onDismiss?: () => void;
}

export const ErrorDisplay = ({error, onRetry, onDismiss}: ErrorDisplayProps) => {
    const [showDetails, setShowDetails] = useState(false);

    return (
        <Box
            sx={{
                maxWidth: 600,
                mx: 'auto',
            }}
            role="alert"
            aria-live="assertive"
        >
            <Alert
                severity="error"
                icon={<ErrorOutlineIcon/>}
                sx={{
                    bgcolor: 'rgba(255, 82, 82, 0.1)',
                    border: '1px solid #ff5252',
                    borderRadius: 0,
                    p: 2,
                }}
            >
                <Typography variant="body1" sx={{mb: error.details ? 1 : 2}}>
                    {error.message}
                </Typography>

                {error.details && (
                    <>
                        <Button
                            size="small"
                            onClick={() => setShowDetails(!showDetails)}
                            sx={{mb: 1, p: 0, minWidth: 0, textTransform: 'none'}}
                        >
                            {showDetails ? 'Hide' : 'Show'} technical details
                        </Button>
                        <Collapse in={showDetails}>
                            <Box
                                sx={{
                                    mt: 1,
                                    p: 1,
                                    bgcolor: 'rgba(0, 0, 0, 0.2)',
                                    fontFamily: 'monospace',
                                    fontSize: 12,
                                    overflowX: 'auto',
                                }}
                            >
                                {error.code && (
                                    <Typography variant="caption" component="div" sx={{fontFamily: 'monospace'}}>
                                        Code: {error.code}
                                    </Typography>
                                )}
                                <Typography variant="caption" component="div" sx={{fontFamily: 'monospace'}}>
                                    {error.details}
                                </Typography>
                            </Box>
                        </Collapse>
                    </>
                )}

                <Box sx={{display: 'flex', gap: 1, mt: 2}}>
                    {onRetry && (
                        <Button
                            variant="contained"
                            size="small"
                            onClick={onRetry}
                            sx={{bgcolor: 'primary.main'}}
                            autoFocus
                        >
                            Try Again
                        </Button>
                    )}
                    {onDismiss && (
                        <Button variant="outlined" size="small" onClick={onDismiss}>
                            Dismiss
                        </Button>
                    )}
                </Box>
            </Alert>
        </Box>
    );
};
