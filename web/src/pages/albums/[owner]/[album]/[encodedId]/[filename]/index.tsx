import App from 'src/App';

// Media page - routing is handled internally in App component
export default function MediaPage() {
  return <App />;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
