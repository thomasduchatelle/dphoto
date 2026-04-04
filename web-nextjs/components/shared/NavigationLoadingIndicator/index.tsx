'use client';

import NextTopLoader from 'nextjs-toploader';

export const NavigationLoadingIndicator = () => {
    return (
        <NextTopLoader
            color="#185986"
            height={3}
            showSpinner={false}
            shadow="0 0 10px #185986, 0 0 5px #185986"
        />
    );
};
