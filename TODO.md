# TODO

## In Progress

- [ ] Implement pagination (cf. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/13)
    - [x] get-members
    - [x] get-members-with-pin
    - [ ] get-groups
    - [ ] get-doors
    - [ ] get-acl
    - [ ] get-acl-with-pin
    - [ ] get-acl-file
    - [ ] get-acl-drive
    - [ ] compare-acl
    - [ ] compare-acl-with-pin
    - [ ] compare-acl-summary
    - [ ] load-acl
    - [ ] load-acl-with-pin
    - [ ] CHANGELOG
    - [ ] README
    - [ ] nightly builds

    - https://gethelp.wildapricot.com/en/articles/2911-preparing-your-api-integrations-for-pagination
    - https://gethelp.wildapricot.com/en/articles/2051-api-update-returned-items-limited-to-100-per-request


## TODO

- [ ] Clean up endOfYear in rules::MakeACL:
```
func (rules *Rules) MakeACL(members types.Members, doors []string) (*ACL, error) {
    ...
    for _, m := range members.Members {
        r := record{
            ...
            ...
            EndDate:   plusOneDay(endOfYear())),
    ...


func (rules *Rules) MakeACLWithPIN(members types.Members, doors []string) (*ACL, error) {
    ...
    for _, m := range members.Members {
        r := record{
            ...
            ...
            EndDate:   plusOneDay(endOfYear())),
    ...
```
- [ ] // FIXME EndDate: 
- [ ] // FIXME use date.Equal
- [ ] // FIXME double check (end date has changed)

- [ ] Use templates for report output
- [ ] Implement generalized struct transcoding

## NOTES
