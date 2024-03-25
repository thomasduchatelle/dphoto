# feature branch

1. semver -> what's the target version ?
2. +update-snapshots -> update the snapshot + commit => ABORT
3. Continuous Integration
   1. test GO
   2. test TS
4. post actions
   1. +next -> deploy on `next` environment
   2. [OPTIONAL] +mr -> create a MR into `main`

# main branch

1. semver -> what's the final version ?
2. Continuous Integration
    1. test GO
    2. test TS
3. Build & Deploy
   1. Deploy -> `next`
   2. Deploy -> `live`
   3. Create GITHUB Release
4. Git TAG
