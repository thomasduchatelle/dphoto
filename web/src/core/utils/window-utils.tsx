import {useEffect, useState} from "react";

export interface WindowDimension {
    width: number
    height: number
}

function getWindowDimensions(): WindowDimension {
    const {innerWidth: width, innerHeight: height} = window;
    return {
        width,
        height
    };
}

export default function useWindowDimensions(): WindowDimension {
    const [windowDimensions, setWindowDimensions] = useState({
        width: 1024,
        height: 720,
    });

    useEffect(() => {
        function handleResize() {
            setWindowDimensions(getWindowDimensions());
        }

        handleResize()
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    return windowDimensions;
}