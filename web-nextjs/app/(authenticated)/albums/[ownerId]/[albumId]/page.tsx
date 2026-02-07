import {Box, Button, Typography} from '@mui/material';
import Link from '@/components/Link';

interface AlbumPageParams {
    ownerId: string;
    albumId: string;
}

export default async function AlbumPage({
                                            params,
                                        }: {
    params: Promise<AlbumPageParams>;
}) {
    const {ownerId, albumId} = await params;

    return (
        <Box sx={{padding: 3}}>
            <Typography variant="h4">Album: {albumId}</Typography>
            <Typography>Owner: {ownerId}</Typography>
            <Typography sx={{marginTop: 2}}>
                Album viewing will be implemented in Epic 2.
            </Typography>
            <Button component={Link} href="/" prefetch={false} sx={{marginTop: 2}}>
                Back to Albums
            </Button>
        </Box>
    );
}
