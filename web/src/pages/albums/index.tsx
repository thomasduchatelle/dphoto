// Albums page - routing is handled internally by React Router in GeneralRouter (via _layout.tsx)
export default function AlbumsPage() {
  return null;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
