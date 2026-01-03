import React from 'react';
import {Story} from '@ladle/react';
import {ErrorDisplay} from '../components/ErrorBoundary/ErrorDisplay';

export default {
    title: 'Layout/ErrorBoundary',
};

type ErrorDisplayProps = React.ComponentProps<typeof ErrorDisplay>;

const ErrorDisplayWrapper: Story<ErrorDisplayProps> = (args) => (
    <ErrorDisplay {...args} />
);

export const Default: Story<ErrorDisplayProps> = (args) => <ErrorDisplayWrapper {...args} />;
Default.args = {
    error: new Error('An unexpected error occurred while processing your request'),
};

export const NetworkError: Story<ErrorDisplayProps> = (args) => <ErrorDisplayWrapper {...args} />;
NetworkError.args = {
    error: new Error('Network request failed: Unable to fetch data from the server'),
};

export const ValidationError: Story<ErrorDisplayProps> = (args) => <ErrorDisplayWrapper {...args} />;
ValidationError.args = {
    error: new Error('Validation failed: Invalid input provided'),
};
