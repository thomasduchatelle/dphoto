# Waku Migration Decision Log

## Migration Topics to Discuss

- Build System & Configuration
- Routing Architecture  
- Server-Side Rendering (SSR)
- React Version Compatibility
- State Management & Data Fetching
- Styling & CSS
- Development Experience
- Testing & Quality Assurance
- Compatibility of libraries with the new version of React
- Visual testing with Storybook or an alternative

## Decisions Made

### Target Runtime Architecture

**Decision**: Use Option 2 without CloudFront - API Gateway + Lambda for SSR

**Architecture**:
