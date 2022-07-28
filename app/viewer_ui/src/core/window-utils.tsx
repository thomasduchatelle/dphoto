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
  const [windowDimensions, setWindowDimensions] = useState(getWindowDimensions());

  useEffect(() => {
    function handleResize() {
      setWindowDimensions(getWindowDimensions());
    }

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return windowDimensions;
}