import App from 'src/App';

// Albums page - routing is handled internally by React Router in App component
export default function AlbumsPage() {
  return <App />;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
