'use client';

import {Box, Button, Collapse, Typography} from '@mui/material';
import {useState} from 'react';

export interface CollapsibleDetailsProps {
    details: string;
}

export const CollapsibleDetails = ({details}: CollapsibleDetailsProps) => {
    const [showDetails, setShowDetails] = useState(false);

    return (
        <Box sx={{width: '100%', maxWidth: 600, textAlign: 'left'}}>
            <Button
                size="small"
                onClick={() => setShowDetails(!showDetails)}
                sx={{mb: 1, p: 0, minWidth: 0, textTransform: 'none', color: 'rgba(255,255,255,0.7)'}}
            >
                {showDetails ? 'Hide' : 'Show'} technical details
            </Button>
            <Collapse in={showDetails}>
                <Box
                    sx={{
                        mt: 1,
                        p: 2,
                        bgcolor: 'rgba(0, 0, 0, 0.3)',
                        fontFamily: 'monospace',
                        fontSize: 12,
                        overflowX: 'auto',
                        borderRadius: 0,
                        color: 'rgba(255,255,255,0.7)',
                    }}
                >
                    <Typography variant="caption" component="pre" sx={{fontFamily: 'monospace', whiteSpace: 'pre-wrap'}}>
                        {details}
                    </Typography>
                </Box>
            </Collapse>
        </Box>
    );
};
