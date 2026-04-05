import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {NoMedia} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const meta = {
    title: 'Catalog/NoMedia',
    parameters: {
        layout: 'fullscreen',
    },
    component: NoMedia,
    decorators: [(Story) => <AppBackground><Story/></AppBackground>],
} satisfies Meta<typeof NoMedia>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
