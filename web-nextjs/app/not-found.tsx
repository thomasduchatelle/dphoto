import {Box, Button, Paper, Typography} from '@mui/material';
import {SearchOff as SearchOffIcon} from '@mui/icons-material';
import Link from '@/components/Link';

export default function NotFound() {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                p: 2,
            }}
        >
            <Paper
                elevation={3}
                sx={{
                    p: 4,
                    maxWidth: 500,
                    textAlign: 'center',
                }}
            >
                <SearchOffIcon
                    sx={{
                        fontSize: 64,
                        color: 'text.secondary',
                        mb: 2,
                    }}
                />
                <Typography variant="h4" gutterBottom>
                    Page Not Found
                </Typography>
                <Typography variant="body1" color="text.secondary" sx={{mb: 3}}>
                    The page you are looking for does not exist
                </Typography>
                <Link href="/">
                    <Button variant="contained" sx={{bgcolor: 'primary.main'}}>
                        Go Home
                    </Button>
                </Link>
            </Paper>
        </Box>
    );
}
