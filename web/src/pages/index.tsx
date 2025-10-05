import App from '../App';

// Root page - all routing is handled internally by React Router in App component
export default function IndexPage() {
  return <App />;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
