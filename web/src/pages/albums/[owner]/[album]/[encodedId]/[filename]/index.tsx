import GeneralRouter from "../../../../../../pages-old/GeneralRouter";

// Media page - routing is handled internally in GeneralRouter (via _layout.tsx)
export default function MediaPage() {
    return <GeneralRouter/>

}

export const getConfig = async () => {
    return {
        render: 'dynamic',
    };
};
