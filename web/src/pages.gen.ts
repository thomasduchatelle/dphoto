// deno-fmt-ignore-file
// biome-ignore format: generated types do not need formatting
// prettier-ignore
import type { PathsForPages, GetConfigResponse } from 'waku/router';

// prettier-ignore
import type { getConfig as File_AlbumsOwnerAlbumEncodedIdFilenameIndex_getConfig } from './pages/albums/[owner]/[album]/[encodedId]/[filename]/index';
// prettier-ignore
import type { getConfig as File_AlbumsOwnerAlbumIndex_getConfig } from './pages/albums/[owner]/[album]/index';
// prettier-ignore
import type { getConfig as File_AlbumsIndex_getConfig } from './pages/albums/index';
// prettier-ignore
import type { getConfig as File_Index_getConfig } from './pages/index';

// prettier-ignore
type Page =
| ({ path: '/albums/[owner]/[album]/[encodedId]/[filename]' } & GetConfigResponse<typeof File_AlbumsOwnerAlbumEncodedIdFilenameIndex_getConfig>)
| ({ path: '/albums/[owner]/[album]' } & GetConfigResponse<typeof File_AlbumsOwnerAlbumIndex_getConfig>)
| ({ path: '/albums' } & GetConfigResponse<typeof File_AlbumsIndex_getConfig>)
| ({ path: '/' } & GetConfigResponse<typeof File_Index_getConfig>);

// prettier-ignore
declare module 'waku/router' {
  interface RouteConfig {
    paths: PathsForPages<Page>;
  }
  interface CreatePagesConfig {
    pages: Page;
  }
}
