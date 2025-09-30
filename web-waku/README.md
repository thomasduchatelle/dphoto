# Web-Waku

This is a minimal Waku-based web application that runs alongside the existing Create React App (CRA) application.

## Architecture

- **Deployment Path**: `/waku` (served via API Gateway + S3)
- **Dev Server Port**: 3001 (to avoid conflicts with CRA on port 3000)
- **Framework**: Waku (React Server Components framework)
- **Routing**: File-based routing (Waku convention)

## Development

### Prerequisites

```bash
npm install
```

### Run Development Server

```bash
npm run dev
```

The development server will start on http://localhost:3001

### Build

```bash
npm run build
```

This creates a production build in the `dist/` directory:
- `dist/public/` - Static assets deployed to S3
- `dist/server/` - Server-side code (for SSR, future enhancement)

## Project Structure

```
web-waku/
├── src/
│   ├── components/     # React components
│   ├── pages/          # File-based routes
│   │   ├── _layout.tsx # Root layout
│   │   ├── index.tsx   # Home page (/waku)
│   │   └── about.tsx   # About page (/waku/about)
│   └── styles.css      # Global styles (Tailwind)
├── public/             # Static assets
├── waku.config.ts      # Waku configuration
└── package.json
```

## Deployment

The application is automatically deployed as part of the CI/CD pipeline:

1. **Build**: `make build-waku` creates the production build
2. **Deploy**: CDK deploys `dist/public/` to S3 bucket
3. **Access**: Available at `https://<domain>/waku`

## Testing

Currently, there are no specific tests configured. The `make test-waku` target runs a build as a smoke test.

## Configuration

### Port Configuration

The dev server port is configured in `waku.config.ts`:

```typescript
export default defineConfig({
  vite: {
    server: {
      port: 3001, // Different from CRA (3000)
    },
  },
});
```

## Integration with Existing App

This Waku app runs in parallel with the existing CRA application:

- **CRA App**: Available at `/` (http://localhost:3000 in dev)
- **Waku App**: Available at `/waku` (http://localhost:3001 in dev)

Both can run simultaneously during development without conflicts.

## Future Enhancements

According to the migration plan (see `specs/2025-09_Waku_Migration.md`):

- Phase 1: Basic deployment (current state)
- Phase 2: SSR with API Gateway + Lambda
- Phase 3+: Component migration and visual testing
