# TODO

## In Progress


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
