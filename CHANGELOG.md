# CHANGELOG

## Unreleased

### Updated
1. Updated to Go 1.24.


## [0.8.10](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.10) - 2025-01-30

### Added
1. ARMv6 build target (RaspberryPi ZeroW).


## [0.8.9](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.9) - 2024-09-06

### Added
1. TCP/IP support.

### Updated
1. Replaced date pointers with concrete types.
2. Replaced the 'start of next year' default card 'valid until' date with end of year.
3. Updated to Go 1.23.


## [0.8.8](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.8) - 2024-03-27

### Updated
1. Bumped Go version to 1.22
2. Reworked member initialisation to resolve group names against system group names using
   group ID because Wild Apricot does not keep the member groups fields consistent.
3. Fixed bug that ignored '0' when normalising strings.


## [0.8.7](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.7) - 2023-12-01

### Changed
1. Maintenance release for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.7


## [0.8.6](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.6) - 2023-08-30

### Updated
1. Replaced os.Rename with lib implementation for tmpfs support.


## [0.8.5](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases/tag/v0.8.5) - 2023-06-13

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


## Older

| *Version* | *Description*                                                                   |
| --------- | ------------------------------------------------------------------------------- |
| v0.7.2    | Maintenance release to update dependencies on `uhppote-core` and `uhppoted-lib` |
| v0.7.1    | Maintenance release to update dependencies on `uhppote-core` and `uhppoted-lib` |
| v0.7.0    | Added support for time profiles from the extended API                           |
| v0.6.12   | Implemented prepended facility code for card numbers retrieved without          |
| v0.6.10   | Initial release                                                                 |
