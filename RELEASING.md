# Release process

Copy these checklists into the release pull request for your convenience:

Prerequisites:
- [ ] Create a milestone for the version
- [ ] Merge all PRs you want to include
- [ ] Mark all solved issues and PRs with the milestone for future reference
- [ ] Check if there is a new LMDB version and update if needed, see `update-lmdb.sh`
- [ ] Create a draft release on Github and let it auto-generate the changelog
- [ ] Add LMDB C library changes to draft release changelog, when applicable

In this pull request:
- [ ] Update CI for newer Go releases
- [ ] Copy the draft release's changelog into CHANGES.md

After merging this pull request:
- [ ] Publish the draft release
- [ ] Close the milestone
