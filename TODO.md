# TODO - v0.7.1

- [x] Update to latest `uhppote-core` and `uhppoted-lib`
- [ ] Fix unit test:
```
go test ./...
--- FAIL: TestAsTable (0.00s)
    acl_test.go:97: Invalid ACL table row
           expected:[1000001 1880-02-29 2022-01-31 Y Y 29 Y]
           got:     [1000001 1880-02-29 2023-01-31 Y Y 29 Y]
    acl_test.go:97: Invalid ACL table row
           expected:[2000001 1981-07-01 2022-01-31 N N N N]
           got:     [2000001 1981-07-01 2023-01-31 N N N N]
    acl_test.go:97: Invalid ACL table row
           expected:[6000001 2021-01-01 2021-06-30 Y N N N]
           got:     [6000001 2022-01-01 2021-06-30 Y N N N]
--- FAIL: TestHash (0.00s)
    acl_test.go:192: Invalid ACL hash - expected:2257f356a9efe68827e7324d4cd68f73b0e9127a1dac65af5659cb87578ef5dc, got:429e2af6ecede056c0eb9bbbd16eb9e0582827434e0e3c3bd9487c66d0d4c192
    acl_test.go:222: Invalid ACL hash - expected:2257f356a9efe68827e7324d4cd68f73b0e9127a1dac65af5659cb87578ef5dc, got:429e2af6ecede056c0eb9bbbd16eb9e0582827434e0e3c3bd9487c66d0d4c192
```

## TODO

- [ ] Use templates for report output
- [ ] Implement generalized struct transcoding

## NOTES