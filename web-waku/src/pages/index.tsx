import {Counter} from "../components/counter";
import {PageProps} from "waku/router";
import {TheRedirector} from "../components/theRedirector";

export default async function HomePage({query}: PageProps<'/'>) {
    const search = new URLSearchParams(query)
    const initialCount = parseInt(search.get("count") ?? "42")

    const data = await getData();

    return (
        <div>
            <title>{data.title}</title>
            <h1 className="text-4xl font-bold tracking-tight">{data.headline}</h1>
            <p>{data.body}</p>
            <Counter startFrom={initialCount}/>
            <TheRedirector count={initialCount}/>
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
        render: 'dynamic',
    } as const;
};
