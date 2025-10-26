import GeneralRouter from "../../pages-old/GeneralRouter";

export default function AlbumsPage() {
    return <GeneralRouter/>;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
