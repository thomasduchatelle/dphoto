import {beforeEach, describe, expect, it, vi} from 'vitest';
import {render, screen} from '@testing-library/react';
import {ErrorBoundary} from './index';
import React, {ReactNode} from 'react';

const ThrowError = ({shouldThrow}: { shouldThrow: boolean }) => {
    if (shouldThrow) {
        throw new Error('Test error message');
    }
    return <div>No error</div>;
};

const TestWrapper = ({children}: { children: ReactNode }) => (
    <>{children}</>
);

describe.skip('ErrorBoundary', () => {
    beforeEach(() => {
        vi.spyOn(console, 'error').mockImplementation(() => {
        });
    });

    it('should render children when no error occurs', () => {
        render(
            <ErrorBoundary>
                <div>Test content</div>
            </ErrorBoundary>
        );

        expect(screen.getByText('Test content')).toBeInTheDocument();
    });

    it('should render error display when child component throws error', () => {
        render(
            <TestWrapper>
                <ErrorBoundary>
                    <ThrowError shouldThrow={true}/>
                </ErrorBoundary>
            </TestWrapper>
        );

        expect(screen.getByText('An error occurred')).toBeInTheDocument();
        expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();
    });

    it('should log error details to console when error is caught', () => {
        const consoleErrorSpy = vi.spyOn(console, 'error');

        render(
            <TestWrapper>
                <ErrorBoundary>
                    <ThrowError shouldThrow={true}/>
                </ErrorBoundary>
            </TestWrapper>
        );

        expect(consoleErrorSpy).toHaveBeenCalled();
        expect(consoleErrorSpy.mock.calls.some(call =>
            call.some(arg => typeof arg === 'string' && arg.includes('Error caught by ErrorBoundary'))
        )).toBe(true);
    });

    it('should display error details notification', () => {
        render(
            <TestWrapper>
                <ErrorBoundary>
                    <ThrowError shouldThrow={true}/>
                </ErrorBoundary>
            </TestWrapper>
        );

        expect(screen.getByText(/Error details have been logged to the console/i)).toBeInTheDocument();
    });
});
