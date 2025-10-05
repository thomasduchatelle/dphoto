import App from 'src/App';

// Album detail page - routing is handled internally in App component
export default function AlbumDetailPage() {
  return <App />;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
