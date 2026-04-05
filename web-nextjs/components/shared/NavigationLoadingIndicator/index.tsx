'use client';

import NextTopLoader from 'nextjs-toploader';
import {useTheme} from '@mui/material/styles';

export const NavigationLoadingIndicator = () => {
    const theme = useTheme();
    const primaryColor = theme.palette.primary.main;

    return (
        <NextTopLoader
            color={primaryColor}
            height={3}
            showSpinner={false}
            shadow={`0 0 10px ${primaryColor}, 0 0 5px ${primaryColor}`}
        />
    );
};
