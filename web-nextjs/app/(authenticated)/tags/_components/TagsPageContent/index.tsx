'use client';

import {useReducer} from 'react';
import {Box, Typography} from '@mui/material';
import StarOutlineIcon from '@mui/icons-material/StarOutline';
import {FAVORITE_TAG, tagPageReducer, TagPageSelection, tagPageSelector, TagPageState} from '@/domains/catalog/tags';
import {ErrorMessage} from '@/components/ErrorMessage';
import {PageMessage} from '@/components/PageMessage';
import {TagMediaGrid} from '../TagMediaGrid';

export interface TagsPageContentProps {
    initialState: TagPageState;
}

export function TagsPageContent({initialState}: TagsPageContentProps) {
    const [state] = useReducer(tagPageReducer, initialState);
    return <TagsPageContentView {...tagPageSelector(state)}/>;
}

function tagTitle(tag: string): string {
    return tag === FAVORITE_TAG ? 'Favorites' : tag;
}

export type TagsPageContentViewProps = TagPageSelection;

export function TagsPageContentView({tag, medias, isEmpty, error}: TagsPageContentViewProps) {
    if (error) {
        return <ErrorMessage error={error} title="Failed to load the tagged medias"/>;
    }

    return (
        <Box>
            <Typography variant="h1" sx={{mb: 3, display: 'flex', alignItems: 'center', gap: 1}}>
                {tag === FAVORITE_TAG && <StarOutlineIcon sx={{fontSize: '0.9em'}}/>}
                {tagTitle(tag)}
            </Typography>

            {isEmpty ? (
                <PageMessage
                    icon={<StarOutlineIcon/>}
                    title="No favourites yet"
                    message="Mark medias as favourite from an album to see them gathered here."
                />
            ) : (
                <TagMediaGrid medias={medias}/>
            )}
        </Box>
    );
}
