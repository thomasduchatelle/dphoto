import {AlbumPageContent} from './_components/AlbumPageContent';
import {Metadata} from 'next';

interface AlbumPageParams {
    ownerId: string;
    albumId: string;
}

export async function generateMetadata({params}: {
    params: Promise<AlbumPageParams>;
}): Promise<Metadata> {
    await params;

    return {
        title: 'DPhoto',
        description: 'Photo management application',
    };
}

export default async function AlbumPage({
                                            params,
                                        }: {
    params: Promise<AlbumPageParams>;
}) {
    await params;

    return <AlbumPageContent/>;
}
