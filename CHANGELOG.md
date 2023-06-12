# CHANGELOG

## [0.8.5](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.5) - 2023-06-14

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.5


## [0.8.4](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.4) - 2023-03-17

### Added
1. `doc.go` package overview documentation.
2. Added `--with-pin` option for card keypad PIN to `get-members`, `get-acl`, `compare-acl`, 
   and `load-acl` commands

### Updated
1. Fixed initial round of _staticcheck_ lint errors and added _staticcheck_ to
   CI build.
2. Added file lock to `get-acl` and `compare-acl` commands with `--lockfile` command line option
   to optionally set the lockfile filepath.


## [0.8.3](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Added ARM64 to release build artifacts

### Changed
1. Moved default `git` branch to `main`, in line with current development practice.
2. Reworked lockfile implementation to use `flock` _syscall_.
3. Removed _zip_ files from release artifacts (no longer necessary)


## [0.8.2](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.2) - 2022-10-14

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.1


## [0.8.1](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.1


## [0.8.0](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.0) - 2022-07-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.0


## [0.7.3](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.7.3) - 2022-06-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.7.3

