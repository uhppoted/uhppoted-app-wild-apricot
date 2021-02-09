# TODO - v0.6.x

### IN PROGRESS

- [ ] Extended dry-run testing
      - [ ] jumpbox
      - [ ] Append reports, (temporarily) deleting before load-acl
      - [x] Log retention
      - [x] `load-acl`
      - [x] wild-apricot.log rotation
      - [x] `get-members`
      - [x] `get-acl`

- [ ] Get member list
      - retry logic
      - `Profile Last Updated` field doesn't seem to be reliable

- [ ] `get-acl`
      - warnings

- [ ] `compare-acl`
      - warnings

- [ ] README

- [ ] ACL
      - [ ] Unit tests for grant/revoke
      - [ ] Variadic grant/revoke e.g. grant("here", "there", 12345)
      - [x] Default start date
      - [x] Default end date
      - [x] Grant/revoke access
      - [x] Export as TSV
      - [x] Door display order
      - [x] Verify `strict` behaviour
      - [x] Rename `record` to `permissions`

- [ ] Commonalise TSV implementation
- [ ] Commonalise MarshalText implementation

- [x] `get-members`
- [x] `load-acl`
- [x] `get-groups`
- [x] `get-doors`
- [x] `get-acl`
- [x] `version` command
- [x] `help` command
- [x] Move `help` to bottom of listed commands in help text
- [x] Get member groups
- [x] Get auth token
- [x] App skeleton

# TODO

- [ ] Use templates for report output
- [ ] Implement generalized struct transcoding

# NOTES