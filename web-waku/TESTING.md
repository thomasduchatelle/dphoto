# Web Waku Testing Guide

This project uses Playwright for visual regression testing against Ladle component stories.

## Running Tests

### Local Development (Quick)
```bash
# Install dependencies
npm install

# Run tests locally
npm test

# Update snapshots (when intentional changes are made)
npm run test:update

# Run tests in headed mode (see browser)
npm run test:headed

# Run tests in UI mode (interactive)
npm run test:ui
```

### Docker-based Testing (CI-consistent)
For consistent results matching GitHub Actions CI:

```bash
# Run tests in Docker (matches CI environment)
npm run test:docker

# Update snapshots in Docker environment
npm run test:docker-update
```

## Test Structure

- **Ladle Stories**: Component stories are defined in `*.stories.tsx` files
- **Playwright Tests**: Automatically generated visual regression tests for all Ladle stories
- **Screenshots**: Baseline screenshots are stored in `src/components/snapshots.spec.ts-snapshots/`

## CI/CD Integration

The project is configured with GitHub Actions that:
1. Installs dependencies
2. Installs Playwright browsers
3. Runs Ladle in the background
4. Executes Playwright visual regression tests
5. Uploads test results and failed screenshots as artifacts

## Adding New Component Tests

1. Create a `ComponentName.stories.tsx` file
2. Export story functions for different component states
3. Run `npm test` - tests are automatically generated
4. Commit the new baseline screenshots

## Troubleshooting

- **Tests failing locally but passing in CI**: Use Docker-based testing (`npm run test:docker`)
- **Screenshot differences**: Check for animations, fonts, or timing issues
- **Ladle not starting**: Ensure port 61000 is available
