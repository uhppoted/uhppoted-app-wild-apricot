# TODO

## In Progress

- [ ] Implement pagination (cf. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/13)
    - [ ] page delay
    - [x] get-members
    - [x] get-members-with-pin
    - [x] ~~get-updated-members~~
    - [x] get-groups
    - [x] ~~get-doors~~
    - [x] get-acl
    - [x] get-acl-with-pin
    - [x] get-acl-file
    - [x] get-acl-drive
    - [x] compare-acl
    - [x] compare-acl-with-pin
    - [x] compare-acl-summary
    - [x] load-acl
    - [x] load-acl-with-pin
    - [x] CHANGELOG
    - [ ] README
    - [ ] release
        - [ ] set MinPageSize, MinPages, etc
        - [ ] uhppoted-lib
        - [ ] uhppoted-app-wild-apricot

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
