// Media page - routing is handled internally in GeneralRouter (via _layout.tsx)
export default function MediaPage() {
  return null;
}

export const getConfig = async () => {
  return {
    render: 'dynamic',
  };
};
