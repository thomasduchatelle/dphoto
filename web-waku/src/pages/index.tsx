export default async function HomePage() {
    const data = await getData();

    return (
        <div>
            <title>{data.title}</title>
            <h1 className="text-4xl font-bold tracking-tight">{data.headline}</h1>
            <p>{data.body}</p>
        </div>
    );
}

const getData = async () => {
    const data = {
        title: 'Waku',
        headline: 'This is the WAKY page.',
        body: "Your base path doesn't work.",
    };

    return data;
};

export const getConfig = async () => {
    return {
        render: 'static',
    } as const;
};
