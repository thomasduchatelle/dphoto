'use client';

import React, {Component, ReactNode} from 'react';
import {ErrorDisplay} from './ErrorDisplay';

interface ErrorBoundaryProps {
    children: ReactNode;
}

interface ErrorBoundaryState {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = {
            hasError: false,
            error: null,
        };
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return {
            hasError: true,
            error,
        };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error('Error caught by ErrorBoundary:', error);
        console.error('Error stack:', error.stack);
        console.error('Component stack:', errorInfo.componentStack);
    }

    render() {
        if (this.state.hasError && this.state.error) {
            return <ErrorDisplay error={this.state.error} />;
        }

        return this.props.children;
    }
}
