# Release process

v2 releases are tagged `v2.X.Y` (the module path is
`github.com/PowerDNS/lmdb-go/v2`). The v1 series continues with `v1.X.Y` tags
from its own branch.

- Check if we need to update the CI for newer Go releases
- Check if there are new LMDB releases and update if needed, per stream:
  `./update-lmdb.sh 09 <version>` and/or `./update-lmdb.sh 10 <version>`.
  The script re-applies the patch series from `lmdb/patches/<stream>/` and
  regenerates the symbol-rename header; if a patch no longer applies,
  re-derive it and update `lmdb/patches/PATCH-STATUS.md`.
- After an LMDB update, run the full validation:
  - `make all` (both engines, C checks, -race)
  - `make coexist-test` (v1+v2 in one binary; needs network)
  - `make bench` on the release commit and its predecessor, compare with
    `benchstat`
  - Include the relevant upstream changelog from `CHANGES.lmdb09.txt` /
    `CHANGES.lmdb10.txt` in `CHANGES.md`
- Create a milestone for the version
- Merge all PRs you want to include
- Mark all solved issues and PRs with the milestone for future reference
- Create the release on Github and let it auto-generate the changelog
- Close the milestone
- Create a new PR for this:
  - Copy-paste the changelog into CHANGES.md
  - Update the links in README.md
