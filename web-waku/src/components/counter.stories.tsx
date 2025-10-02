import {Counter} from "./counter";
export default {
    meta: {
        width: "large",
        iframed: true
    }
}

export const Zero = () => <Counter></Counter>
Zero.meta = {iframed: true}
export const One = () => <h1>Hi !</h1>
One.meta = {iframed: true}