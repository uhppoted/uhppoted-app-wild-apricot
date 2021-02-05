# TODO - v0.6.x

### IN PROGRESS

- [ ] Extended dry-run testing
      - [ ] wild-apricot.log rotation
      - [ ] `load-acl`
      - [x] `get-members`
      - [x] `get-acl`
      - [ ] jumpbox

- [ ] `load-acl`
       - log
       - report

- [ ] Get member list
      - Implement async fetch
      - Extend timeout (or retry?)

- [ ] ACL
      - [x] Default start date
      - [x] Default end date
      - [x] Grant/revoke access
      - [x] Export as TSV
      - [x] Door display order
      - [ ] Unit tests for grant/revoke
      - [x] Verify `strict` behaviour
      - [ ] Rename `record` to `permissions`

- [ ] Commonalise TSV implementation
- [ ] Commonalise MarshalText implementation

- [x] `get-members`
- [x] `get-groups`
- [x] `get-doors`
- [x] `get-acl`
- [x] `compare-acl`
- [x] `get-acl`
- [x] Get member groups
- [x] App skeleton
- [x] `version` command
- [x] `help` command
- [x] Get auth token
- [x] Move `help` to bottom of listed commands in help text

# TODO

- [ ] Use templates for report output
- [ ] Implement generalized struct transcoding

# NOTES