import {AlbumPageContent} from './_components/AlbumPageContent';
import {Metadata} from 'next';
import {newServerSideRestCatalogAdapter} from '@/domains/catalog/adapters/server-adapter-factory';

interface AlbumPageParams {
    ownerId: string;
    albumId: string;
}

export async function generateMetadata({params}: {
    params: Promise<AlbumPageParams>;
}): Promise<Metadata> {
    const {ownerId, albumId} = await params;

    try {
        const adapter = newServerSideRestCatalogAdapter();
        const albums = await adapter.fetchAlbums();
        const album = albums.find(
            a => a.albumId.owner === ownerId && a.albumId.folderName === albumId
        );

        if (album) {
            return {
                title: `${album.name} - DPhoto`,
                description: `View photos from ${album.name}`,
            };
        }
    } catch (error) {
        console.error('Error fetching album for metadata:', error);
    }

    return {
        title: 'Album - DPhoto',
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
