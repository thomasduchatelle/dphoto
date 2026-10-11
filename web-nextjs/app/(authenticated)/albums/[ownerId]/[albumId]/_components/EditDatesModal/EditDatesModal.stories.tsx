import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {EditDatesModal} from './index';

const meta = {
    title: 'Catalog/EditDatesModal',
    component: EditDatesModal,
    parameters: {
        layout: 'fullscreen',
    },
    args: {
        isOpen: true,
        albumName: 'First Contact',
        startDate: new Date('2063-04-05T00:00:00.000Z'),
        endDate: new Date('2063-04-05T23:59:00.000Z'),
        startAtDayStart: true,
        endAtDayEnd: true,
        isLoading: false,
        isSaveEnabled: true,
        onClose: fn(),
        onSave: fn(),
        onStartDateChange: fn(),
        onEndDateChange: fn(),
        onStartAtDayStartChange: fn(),
        onEndAtDayEndChange: fn(),
    },
} satisfies Meta<typeof EditDatesModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithSpecificTimes: Story = {
    args: {
        albumName: 'Mission to Mars',
        startDate: new Date('2063-04-05T10:30:00.000Z'),
        endDate: new Date('2063-04-07T15:45:00.000Z'),
        startAtDayStart: false,
        endAtDayEnd: false,
    },
};

export const InvalidRange: Story = {
    args: {
        startDate: new Date('2063-04-07T00:00:00.000Z'),
        endDate: new Date('2063-04-05T00:00:00.000Z'),
        dateRangeError: 'The end date cannot be before the start date',
        isSaveEnabled: false,
    },
};

export const OrphanWarning: Story = {
    args: {
        errorCode: 'The dates cannot be updated if it makes some pictures orphaned. Create another album first.',
    },
};

export const SubmissionError: Story = {
    args: {
        errorCode: 'Album dates couldn\u2019t be saved. Refresh your page and retry, or let the developer known.',
    },
};

export const Loading: Story = {
    args: {
        isLoading: true,
        isSaveEnabled: false,
    },
};
