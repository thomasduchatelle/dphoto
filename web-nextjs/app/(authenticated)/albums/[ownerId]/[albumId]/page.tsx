import {AlbumPageContent} from './_components/AlbumPageContent';

interface AlbumPageParams {
    ownerId: string;
    albumId: string;
}

export default async function AlbumPage({
                                            params,
                                        }: {
    params: Promise<AlbumPageParams>;
}) {
    await params;

    return <AlbumPageContent/>;
}
