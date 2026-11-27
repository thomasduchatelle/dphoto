import {Metadata} from 'next';
import {serverSideThunk} from '@/libs/dthunks/server';
import {FAVORITE_TAG, initialTagPageState, tagPageThunks} from '@/domains/catalog/tags';
import {newTagsAdapter} from '@/domains/catalog/tags/adapters/tags-adapter-factory';
import {TagsPageContent} from './_components/TagsPageContent';
import {mustBeAuthenticated} from '../auth-utils';

export const metadata: Metadata = {
    title: 'Tags - DPhoto',
    description: 'Browse your medias by tag',
};

export default async function TagsPage() {
    await mustBeAuthenticated();

    const loadTagMedias = serverSideThunk(
        tagPageThunks.loadTagMedias,
        {adapter: newTagsAdapter()},
    );
    const tagPageState = await loadTagMedias(initialTagPageState(FAVORITE_TAG), FAVORITE_TAG);

    return <TagsPageContent initialState={tagPageState}/>;
}
