import Link from '@/components/Link';
import {Box, Button, Paper, Typography} from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';

export default function LogoutPage() {

    return (
        <Box
            sx={{
                display: 'flex',
                minHeight: '100vh',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            <Paper
                sx={{
                    mx: 2,
                    maxWidth: '28rem',
                    p: 4,
                    textAlign: 'center',
                }}
            >
                <Box sx={{mb: 3, display: 'flex', justifyContent: 'center'}}>
                    <CheckCircleOutlineIcon
                        sx={{
                            fontSize: '4rem',
                            color: 'success.main',
                        }}
                    />
                </Box>
                <Typography
                    variant="h4"
                    sx={{
                        mb: 2,
                        fontWeight: 'bold',
                    }}
                >
                    Successfully Logged Out
                </Typography>
                <Typography
                    variant="body1"
                    sx={{
                        mb: 4,
                        color: 'text.secondary',
                    }}
                >
                    You have been successfully logged out of your account.
                </Typography>
                <Box sx={{display: 'flex', justifyContent: 'center'}}>
                    <Button
                        component={Link}
                        href="/"
                        prefetch={false}
                        variant="contained"
                        sx={{
                            borderRadius: '24px',
                            px: 3,
                            py: 1.5,
                        }}
                    >
                        Sign In Again
                    </Button>
                </Box>
            </Paper>
        </Box>
    );
}
