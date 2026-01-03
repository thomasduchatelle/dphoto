import GeneralRouter from "../pages-old/GeneralRouter";

// Root page - all routing is handled internally by React Router in GeneralRouter (via _layout.tsx)
export default function IndexPage() {
    return <GeneralRouter/>
}

export const getConfig = async () => {
    return {
        render: 'dynamic',
    };
};
