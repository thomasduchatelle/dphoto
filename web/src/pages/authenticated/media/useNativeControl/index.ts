import {useEffect} from "react";
import {SwipeableHandlers, useSwipeable} from "react-swipeable";

export enum Key {
  Left = "Left",
  Right = "Right",
  D = "D",
  Del = "Del",
  Esc = "Esc",
}

// Listen for device native input (keyboard, swipe, ...)
export function useNativeControl(callback: (key: Key) => void, ...keys: Key[]): SwipeableHandlers {
  const handlers = useSwipeable({
    onSwiped: (eventData) => {
      const key = swipeToKey(eventData.dir)
      if (key && keys.indexOf(key) >= 0) {
        callback(key)
      }
    },
  });

  useEffect(() => {
    const handler = (evt: KeyboardEvent) => {
      const key = keyboardToKey(evt.key)
      if (key && keys.indexOf(key) >= 0) {
        evt.preventDefault()
        callback(key)
      }
    }

    window.addEventListener("keydown", handler);

    return () => {
      window.removeEventListener("keydown", handler);
    };
  }, [callback, keys]);

  return handlers
}

function keyboardToKey(key: string): Key | null {
  switch (key) {
    case "d":
    case "D":
      return Key.D

    case "Backspace":
      return Key.Del

    case "ArrowLeft":
      return Key.Left

    case "ArrowRight":
      return Key.Right

    case "Escape":
      return Key.Esc

    default:
      return null
  }
}

function swipeToKey(swipeDirection: string): Key | null {
  switch (swipeDirection) {
    case "Right":
      return Key.Left

    case "Left":
      return Key.Right

    default:
      return null
  }
}
