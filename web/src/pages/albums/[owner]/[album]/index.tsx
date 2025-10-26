import GeneralRouter from "../../../../pages-old/GeneralRouter";

// Album detail page - routing is handled internally in GeneralRouter (via _layout.tsx)
export default function AlbumDetailPage() {
    return <GeneralRouter/>

}

export const getConfig = async () => {
    return {
        render: 'dynamic',
    };
};
